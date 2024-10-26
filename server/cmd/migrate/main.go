/*
Credits: https://github.com/jatindotdev/tinybits/blob/main/migrations/migrate.go
*/

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
)

var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
	dir   = flags.String("dir", "./migrations", "directory with migration files")
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func main() {
	if err := flags.Parse(os.Args[1:]); err != nil {
		log.Fatalf("Error parsing flags: %v", err)
	}
	args := flags.Args()

	if len(args) < 1 {
		flags.Usage()
		return
	}

	command := args[0]
	dbString := os.Getenv("DATABASE_URL")
	log.Println(dbString)

	if dbString == "" {
		fmt.Println("goose: missing DATABASE_URL environment variable")
		return
	}

	db, err := goose.OpenDBWithDriver("postgres", dbString)
	if err != nil {
		fmt.Printf("goose: failed to open DB: %v\n", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			fmt.Printf("goose: failed to close DB: %v\n", err)
		}
	}()

	arguments := []string{}
	if len(args) > 1 {
		arguments = append(arguments, args[1:]...)
	}

	if err := goose.RunContext(context.Background(), command, db, *dir, arguments...); err != nil {
		fmt.Printf("goose %v: %v\n", command, err)
	}
}
