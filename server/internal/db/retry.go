package db

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
)

const (
	maxRetries              = 3
	retryInterval           = 1500 * time.Millisecond
	maxQueueSize            = 1000
	queueTimeout            = 5 * time.Second
	healthCheckPeriod       = 30 * time.Second
	circuitTickerInterval   = 5 * time.Minute
	circuitBreakerThreshold = 10
	circuitBreakerResetTime = 1 * time.Minute
)

type ConnectionState int32

const (
	StateHealthy ConnectionState = iota
	StateReconnecting
	StateFailed
	StateCircuitOpen
)

type queuedOperation struct {
	operation func() error
	resultCh  chan error
	ctx       context.Context
}

type ConnectionManager struct {
	mu              sync.RWMutex
	state           atomic.Int32
	operationQueue  chan *queuedOperation
	quit            chan struct{}
	failureCount    atomic.Int32
	lastFailureTime atomic.Int64
	metrics         *ConnectionMetrics
}

type ConnectionMetrics struct {
	totalOperations  atomic.Int64
	failedOperations atomic.Int64
	queuedOperations atomic.Int64
	avgResponseTime  atomic.Int64
	circuitBreaks    atomic.Int64
}

func NewConnectionManager() *ConnectionManager {
	cm := &ConnectionManager{
		operationQueue: make(chan *queuedOperation, maxQueueSize),
		quit:           make(chan struct{}),
		metrics:        &ConnectionMetrics{},
	}
	cm.state.Store(int32(StateHealthy))

	go cm.processQueue()
	go cm.monitorMetrics()

	return cm
}

func (cm *ConnectionManager) RetryWithSingleFlight(ctx context.Context, operation func() error) error {
	start := time.Now()
	defer func() {
		cm.metrics.totalOperations.Add(1)
		cm.metrics.avgResponseTime.Store(time.Since(start).Milliseconds())
	}()

	if cm.isCircuitOpen() {
		if !cm.shouldAttemptReset() {
			return ErrCircuitOpen
		}
	}

	if err := cm.executeWithContext(ctx, operation); err == nil {
		cm.setState(StateHealthy)
		return nil
	}

	cm.failureCount.Add(1)
	cm.lastFailureTime.Store(time.Now().Unix())
	cm.metrics.failedOperations.Add(1)

	if cm.failureCount.Load() >= circuitBreakerThreshold {
		cm.setState(StateCircuitOpen)
		cm.metrics.circuitBreaks.Add(1)
		return ErrCircuitOpen
	}

	currentState := ConnectionState(cm.state.Load())
	if currentState == StateHealthy {
		cm.setState(StateReconnecting)
		return cm.performRecovery(ctx, operation)
	}

	if currentState == StateReconnecting {
		return cm.queueOperation(ctx, operation)
	}

	return ErrRecoveryFailed
}

func (cm *ConnectionManager) executeWithContext(ctx context.Context, operation func() error) error {
	done := make(chan error, 1)

	go func() {
		done <- operation()
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (cm *ConnectionManager) queueOperation(ctx context.Context, operation func() error) error {
	cm.metrics.queuedOperations.Add(1)

	qOp := &queuedOperation{
		operation: operation,
		resultCh:  make(chan error, 1),
		ctx:       ctx,
	}

	select {
	case cm.operationQueue <- qOp:
		select {
		case result := <-qOp.resultCh:
			return result
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(queueTimeout):
			return ErrQueueTimeout
		}
	default:
		return ErrQueueFull
	}
}

func (cm *ConnectionManager) processQueue() {
	for {
		select {
		case <-cm.quit:
			return
		case qOp := <-cm.operationQueue:
			select {
			case <-qOp.ctx.Done():
				qOp.resultCh <- qOp.ctx.Err()
				continue
			default:
			}

			if ConnectionState(cm.state.Load()) == StateHealthy {
				err := qOp.operation()
				qOp.resultCh <- err
			} else {
				qOp.resultCh <- ErrNotHealthy
			}
		}
	}
}

func (cm *ConnectionManager) monitorMetrics() {
	ticker := time.NewTicker(circuitTickerInterval)
	defer ticker.Stop()

	for {
		select {
		case <-cm.quit:
			return
		case <-ticker.C:
			slog.Info("Connection metrics",
				"total_operations", cm.metrics.totalOperations.Load(),
				"failed_operations", cm.metrics.failedOperations.Load(),
				"queued_operations", cm.metrics.queuedOperations.Load(),
				"avg_response_time_ms", cm.metrics.avgResponseTime.Load(),
				"circuit_breaks", cm.metrics.circuitBreaks.Load(),
			)
		}
	}
}

func (cm *ConnectionManager) isCircuitOpen() bool {
	return ConnectionState(cm.state.Load()) == StateCircuitOpen
}

func (cm *ConnectionManager) shouldAttemptReset() bool {
	lastFailure := time.Unix(cm.lastFailureTime.Load(), 0)
	return time.Since(lastFailure) > circuitBreakerResetTime
}

func (cm *ConnectionManager) setState(state ConnectionState) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.state.Store(int32(state))
	slog.Info("Connection state changed", "state", state)
}

func (cm *ConnectionManager) GetMetrics() ConnectionMetrics {
	return ConnectionMetrics{
		totalOperations:  atomic.Int64{},
		failedOperations: atomic.Int64{},
		queuedOperations: atomic.Int64{},
		avgResponseTime:  atomic.Int64{},
		circuitBreaks:    atomic.Int64{},
	}
}

func (cm *ConnectionManager) Close() {
	close(cm.quit)
}

var (
	ErrReconnecting   = errors.New("connection recovery in progress")
	ErrRecoveryFailed = errors.New("connection recovery failed")
	ErrUnknownState   = errors.New("unknown connection state")
	ErrQueueFull      = errors.New("operation queue is full")
	ErrQueueTimeout   = errors.New("operation queue timeout")
	ErrNotHealthy     = errors.New("connection not healthy")
	ErrCircuitOpen    = errors.New("circuit breaker is open")
)

func (cm *ConnectionManager) performRecovery(ctx context.Context, operation func() error) error {
	retryCount := 0

	for retryCount < maxRetries {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(retryInterval):
		}

		err := operation()
		if err == nil {
			cm.failureCount.Store(0)
			cm.setState(StateHealthy)
			return nil
		}

		retryCount++
		slog.Warn("Retry attempt failed",
			"attempt", retryCount,
			"maxRetries", maxRetries,
			"error", err)
	}

	cm.setState(StateFailed)
	return ErrRecoveryFailed
}
