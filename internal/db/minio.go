package db

import (
	"context"
	"log"
	"log/slog"

	"github.com/HarshPatel5940/CodeFlick/internal/utils"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func CreateMinioClient() *minio.Client {
	var SSLPolicy bool
	if utils.GetEnv("MINIO_SSL_POLICY", "false") == "false" {
		SSLPolicy = false
	} else {
		SSLPolicy = true
	}
	minioClient, err := minio.New(utils.GetEnv("MINIO_ENDPOINT"), &minio.Options{
		Creds: credentials.NewStaticV4(
			utils.GetEnv("MINIO_ACCESS_KEY"),
			utils.GetEnv("MINIO_ACCESS_SECRET"),
			""),
		Secure: SSLPolicy,
	})

	if err != nil {
		panic(err)
	}
	return minioClient
}

func InitMinioClient(MinioClient *minio.Client) {
	bucketName := utils.GetEnv("MINIO_BUCKET_NAME", "codeflick")
	bucketExists, err := MinioClient.BucketExists(context.Background(), bucketName)
	if err != nil {
		slog.Error("Error while checking if bucket exists")
		log.Panic(err)
	}

	if !bucketExists {
		err = MinioClient.MakeBucket(
			context.Background(),
			bucketName,
			minio.MakeBucketOptions{
				Region:        "ap-south-1",
				ObjectLocking: false,
			},
		)

		if err != nil {
			slog.Error("Error while creating bucket")
			log.Panic(err)
			return
		}
		bucketExists = true
		slog.Info("Successfully created mybucket.")
	}

	if bucketExists {
		err = MinioClient.SetBucketPolicy(
			context.Background(),
			bucketName,
			utils.GetBucketPolicy(bucketName),
		)

		if err != nil {
			slog.Error("Error while setting bucket policy")
			log.Panic(err)
			return
		}

		slog.Info("Successfully set bucket policy.")
	}

	slog.Info("Bucket Initialized Successfully.")
}
