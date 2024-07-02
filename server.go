package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

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
	minioBucketName := "codeflick"
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

	e.GET("/storage/buckets", func(c echo.Context) error {
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

	e.POST("/storage/new", func(c echo.Context) error {
		user := c.FormValue("user")
		fileName := c.FormValue("name")
		// email := c.FormValue("email")

		file, err := c.FormFile("file")
		if err != nil {
			return err
		}
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		location := fmt.Sprintf("%s/%s", user, fileName)

		//----------- send to minio docker img
		fileInfo, err := minioClient.PutObject(context.Background(),
			minioBucketName,
			location,
			src,
			file.Size,
			minio.PutObjectOptions{ContentType: file.Header.Get("Content-Type"), Expires: <-time.After(time.Hour * 24 * 30)})

		if err != nil {
			log.Fatalln(err)
		}

		slog.Info("Successfully uploaded", fileInfo.Size)
		return c.String(200, location)
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

/* bucket policy
{
    "Statement": [{
        "Action": ["s3:GetBucketLocation", "s3:ListBucket"],
        "Effect": "Allow",
        "Principal": {
            "AWS": ["*"]
        },
        "Resource": ["arn:aws:s3:::job-offers"]
    }, {
        "Action": ["s3:GetObject"],
        "Effect": "Allow",
        "Principal": {
            "AWS": ["*"]
        },
        "Resource": ["arn:aws:s3:::job-offers/*"]
    }],
    "Version": "2012-10-17"
}

*/
