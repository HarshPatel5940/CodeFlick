package db

import (
	"github.com/HarshPatel5940/CodeFlick/internal/utils"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func CreatePostgresConnection() *sqlx.DB {
	db, err := sqlx.Open("postgres", utils.GetEnv("DATABASE_URL"))

	if err != nil {
		panic(err)
	}

	return db
}
