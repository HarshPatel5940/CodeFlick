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

type MinioHandler struct {
	client *minio.Client
	cm     *ConnectionManager
}

var (
	MinioBucketName string = utils.GetEnv("MINIO_BUCKET_NAME", "codeflick")
	MinioSSLPolicy  string = utils.GetEnv("MINIO_SSL_POLICY", "false")
)

func CreateMinioClient(cm *ConnectionManager) *MinioHandler {
	var client *minio.Client
	var err error
	SSLPolicy := MinioSSLPolicy == "true"

	for attempt := 0; attempt <= maxRetries; attempt++ {
		client, err = minio.New(utils.GetEnv("MINIO_ENDPOINT"), &minio.Options{
			Creds: credentials.NewStaticV4(
				utils.GetEnv("MINIO_ACCESS_KEY"),
				utils.GetEnv("MINIO_ACCESS_SECRET"),
				""),
			Secure: SSLPolicy,
		})

		if err == nil {
			_, err = client.ListBuckets(context.Background())
		}

		if handleError(err, attempt, "initialize minio client") {
			continue
		}

		break
	}

	return &MinioHandler{client: client, cm: cm}
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
		slog.Info("Successfully created the bucket.")
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

func (m *MinioHandler) GetBucketName() string {
	return MinioBucketName
}

func (m *MinioHandler) PutObject(ctx context.Context, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (minio.UploadInfo, error) {
	var info minio.UploadInfo
	err := m.cm.RetryWithSingleFlight(ctx, func() error {
		var err error
		info, err = m.client.PutObject(ctx, MinioBucketName, objectName, reader, objectSize, opts)
		return err
	})
	return info, err
}

func (m *MinioHandler) GetObject(ctx context.Context, objectName string, opts minio.GetObjectOptions) (*minio.Object, error) {
	var obj *minio.Object
	err := m.cm.RetryWithSingleFlight(ctx, func() error {
		var err error
		obj, err = m.client.GetObject(ctx, MinioBucketName, objectName, opts)
		return err
	})
	return obj, err
}

func (m *MinioHandler) RemoveObject(ctx context.Context, objectName string, opts minio.RemoveObjectOptions) error {
	return m.cm.RetryWithSingleFlight(ctx, func() error {
		return m.client.RemoveObject(ctx, MinioBucketName, objectName, opts)
	})
}

func (m *MinioHandler) ListBuckets(ctx context.Context) ([]minio.BucketInfo, error) {
	var buckets []minio.BucketInfo
	err := m.cm.RetryWithSingleFlight(ctx, func() error {
		var err error
		buckets, err = m.client.ListBuckets(ctx)
		return err
	})
	return buckets, err
}

func (m *MinioHandler) ListObjects(ctx context.Context, prefix string, recursive bool) <-chan minio.ObjectInfo {
	return m.client.ListObjects(ctx, MinioBucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: recursive,
	})
}
