package db

import (
	"context"
	"log"
	"log/slog"
	"time"

	"github.com/HarshPatel5940/CodeFlick/internal/utils"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func CreatePostgresConnection(cm *RetryManager) *sqlx.DB {
	var db *sqlx.DB
	var err error

	if err = cm.RetryWithSingleFlight(context.Background(), func() error {
		db, err = sqlx.Open("postgres", utils.GetEnv("DATABASE_URL"))
		if err != nil {
			return err
		}

		if err = db.Ping(); err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Fatal("Error connecting to database", "error", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	go InitPostgresDB(db)

	return db
}

func InitPostgresDB(db *sqlx.DB) {
	if err := goose.SetDialect("postgres"); err != nil {
		slog.Error("Failed to set dialect", "error", err)
		return
	}
	version, err := goose.GetDBVersion(db.DB)
	if err != nil {
		slog.Error("Failed to get database version", "error", err)
		return
	}

	slog.Info("Current migration version", "version", version)
}
