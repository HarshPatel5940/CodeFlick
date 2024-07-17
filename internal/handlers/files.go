package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/HarshPatel5940/CodeFlick/internal/utils"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
)

var MinioBucketName string = utils.GetEnv("MINIO_BUCKET_NAME", "codeflick")

type FileStorageHandler struct {
	db    *sqlx.DB
	minio *minio.Client
}

func NewFilesHandler(db *sqlx.DB, minio *minio.Client) *FileStorageHandler {
	return &FileStorageHandler{db, minio}
}

func (fh FileStorageHandler) UploadFile(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed to get session!",
			"details": err.Error(),
		})
	}
	email := sess.Values["email"].(string)

	fileName := c.FormValue("fileName")
	file, err := c.FormFile("file")

	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	location := fmt.Sprintf("%s/%s", email, fileName)

	fileInfo, err := fh.minio.PutObject(context.Background(),
		MinioBucketName,
		location,
		src,
		file.Size,
		minio.PutObjectOptions{ContentType: file.Header.Get("Content-Type")})
	// 	UserMetadata: map[string]string{         <- Will this leak user mail along with the file?
	// 	"email": email,
	// }
	// Expires: <-time.After(time.Hour * 24 * 30) <- I thought of keeping expiring objects but change of plans i guess

	if err != nil {
		log.Fatalln(err)
	}

	return c.JSON(http.StatusCreated, map[string]any{
		"key":          fileInfo.Key,
		"fileSize":     fileInfo.Size,
		"fileLocation": fileInfo.Location,
		"fileTag":      fileInfo.ETag,
	})
}

func (fh FileStorageHandler) ListBuckets(c echo.Context) error {
	list, err := fh.minio.ListBuckets(context.Background())
	if err != nil {
		resp := minio.ToErrorResponse(err)

		if resp.Code != "" {
			return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
				"success": false,
				"Message": resp.Message,
				"Details": resp,
			})
		}
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
			"success": false,
			"Message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": list,
	})
}
