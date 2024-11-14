package db

import (
	"log/slog"
	"os/exec"
	"time"

	"github.com/HarshPatel5940/CodeFlick/internal/utils"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func CreatePostgresConnection(cm *RetryManager) *sqlx.DB {
	var db *sqlx.DB
	var err error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		db, err = sqlx.Open("postgres", utils.GetEnv("DATABASE_URL"))
		if handleError(err, attempt, "connect to postgres") {
			continue
		}

		err = db.Ping()
		if handleError(err, attempt, "ping postgres") {
			continue
		}

		break
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	go InitPostgresDB()

	return db
}

func handleError(err error, attempt int, operation string) bool {
	if err != nil {
		if attempt < maxRetries {
			slog.Error("Failed to "+operation+", retrying...",
				"attempt", attempt+1,
				"error", err)
			time.Sleep(retryInterval)
			return true
		}
		slog.Error("Failed to "+operation+" after all retries", "error", err)
		panic(err)
	}
	return false
}

func InitPostgresDB() {
	cmd := exec.Command("make", "migrate-status")
	output, err := cmd.CombinedOutput()
	if err != nil {
		slog.Error("Error executing migration command", "error", err)
	}

	if len(output) > 0 {
		slog.Info(string(output))
	}
}
