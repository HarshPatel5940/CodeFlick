package db

import (
	"context"
	"io"
	"log"
	"log/slog"

	"github.com/HarshPatel5940/CodeFlick/internal/utils"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioBucketName string = utils.GetEnv("MINIO_BUCKET_NAME", "codeflick")

type MinioHandler struct {
	client *minio.Client
}

func CreateMinioClient() *MinioHandler {
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
	return &MinioHandler{client: minioClient}
}

func (m *MinioHandler) InitMinioClient() {
	bucketExists, err := m.client.BucketExists(context.Background(), MinioBucketName)
	if err != nil {
		slog.Error("Error while checking if bucket exists")
		log.Panic(err)
	}

	if !bucketExists {
		err = m.client.MakeBucket(
			context.Background(),
			MinioBucketName,
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
		err = m.client.SetBucketPolicy(
			context.Background(),
			MinioBucketName,
			utils.GetBucketPolicy(MinioBucketName),
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

func (m *MinioHandler) PutObject(ctx context.Context, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (minio.UploadInfo, error) {
	return m.client.PutObject(ctx, MinioBucketName, objectName, reader, objectSize, opts)
}

func (m *MinioHandler) GetObject(ctx context.Context, objectName string, opts minio.GetObjectOptions) (*minio.Object, error) {
	return m.client.GetObject(ctx, MinioBucketName, objectName, opts)
}

func (m *MinioHandler) RemoveObject(ctx context.Context, objectName string, opts minio.RemoveObjectOptions) error {
	return m.client.RemoveObject(ctx, MinioBucketName, objectName, opts)
}

func (m *MinioHandler) ListBuckets(ctx context.Context) ([]minio.BucketInfo, error) {
	return m.client.ListBuckets(ctx)
}
