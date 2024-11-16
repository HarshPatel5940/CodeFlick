package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var env = map[string]interface{}{
	"PORT":                    "8080",
	"ENV":                     "development",
	"CLIENT_URL":              "http://localhost:3000",
	"DATABASE_URL":            nil,
	"GOOGLE_CLIENT_ID":        nil,
	"GOOGLE_CLIENT_SECRET":    nil,
	"GOOGLE_CALLBACK_URL":     "http://localhost:8080/auth/google/callback",
	"MINIO_ACCESS_KEY":        nil,
	"MINIO_ACCESS_SECRET":     nil,
	"MINIO_SSL_POLICY":        "false",
	"MINIO_BUCKET_NAME":       "codeflick",
	"MINIO_ENDPOINT":          "localhost:9000",
	"GORILLA_SESSIONS_MAXAGE": "604800",
	"GORILLA_SESSIONS_KEY":    "NotSoSecretKey-ChangeMe-Please",
}

func GetEnv(key string, fallback ...string) string {
	value, exists := os.LookupEnv(key)

	if !exists {
		if len(fallback) > 0 {
			return fallback[0]
		}

		if env[key] != nil {
			return fmt.Sprintf("%v", env[key])
		}

		log.Panicf("Environment variable %s not set", key)
	}

	return value
}

func GetServerAddress() string {
	port := GetEnv("PORT", "8080")
	env := GetEnv("ENV", "development")

	switch env {
	case "production", "prod", "staging", "docker":
		return fmt.Sprintf(":%s", port)
	default:
		return fmt.Sprintf("localhost:%s", port)
	}
}

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	for key, value := range env {
		if value, ok := value.(string); ok {
			GetEnv(key, value)
		} else {
			GetEnv(key)
		}
	}
}
