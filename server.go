package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	Key    = "HelloWorld"
	MaxAge = 86400 * 30
	IsProd = false
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	minioAccessKey := os.Getenv("MINIO_ACCESS_KEY")
	minioAccessSecret := os.Getenv("MINIO_ACCESS_SECRET")
	endpoint := "localhost:9000"
	useSSL := false

	store := sessions.NewCookieStore([]byte(Key))

	store.Options.Path = "/"
	store.Options.MaxAge = MaxAge
	store.Options.HttpOnly = true
	store.Options.Secure = IsProd

	gothic.Store = store

	goth.UseProviders(
		google.New(googleClientID, googleClientSecret, "http://localhost:8080/auth/google/callback", "email", "profile"),
	)

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioAccessKey, minioAccessSecret, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
		// minioClient is now setup
	})

	e.GET("/minio", func(c echo.Context) error {
		list, err := minioClient.ListBuckets(context.Background())
		if err != nil {
			log.Fatalln(err)
		}

		var bucketNames string
		for _, bucket := range list {
			bucketNames += bucket.Name + ", "
		}

		return c.String(200, bucketNames)
	})

	e.GET("/auth/:provider/login", func(c echo.Context) error {
		// try to get the user without re-authenticating
		ctx := context.WithValue(c.Request().Context(), "provider", "google")
		gothUser, err := gothic.CompleteUserAuth(c.Response(), c.Request().WithContext(ctx))
		if err != nil {
			gothic.BeginAuthHandler(c.Response(), c.Request().WithContext(ctx))
		} else {
			log.Println(gothUser)
		}

		return nil
	})

	e.GET("/auth/:provider/callback", func(c echo.Context) error {

		user, err := gothic.CompleteUserAuth(c.Response(), c.Request())
		if err != nil {
			return err
		}
		// store the user in a session

		return c.JSON(200, user)
	})

	slog.Info("Server started at http://127.0.0.1:8080")
	e.Logger.Fatal(e.Start(":8080"))

}
