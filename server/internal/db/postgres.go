package db

import (
	"context"
	"log"
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

	go InitPostgresDB()

	return db
}

func InitPostgresDB() {
	cmd := exec.Command("make", "migrate-status")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Panic("Error executing migration command", "error", err)
	}

	if len(output) > 0 {
		slog.Info(string(output))
	}
}
