package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/HarshPatel5940/CodeFlick/internal/models"
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

	var Gist models.Gist
	Gist.UserID = sess.Values["user_id"].(string)

	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	Gist.GistTitle = c.FormValue("gist_title")
	if Gist.GistTitle == "" {
		Gist.GistTitle = fmt.Sprintf("%d-%s", time.Now().Unix(), file.Filename)
	}

	Gist.ShortUrl = c.FormValue("custom_url")
	if Gist.ShortUrl == "" {
		Gist.ShortUrl = file.Filename
	}

	Gist.IsPublic, err = strconv.ParseBool(c.FormValue("is_public"))
	if err != nil {
		Gist.IsPublic = false
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	fileNameParts := strings.SplitAfter(file.Filename, ".")

	location := fmt.Sprintf("%s/%d.%s",
		Gist.UserID, time.Now().Unix(), fileNameParts[len(fileNameParts)-1],
	)

	fileInfo, err := fh.minio.PutObject(context.Background(),
		MinioBucketName,
		location,
		src,
		file.Size,
		minio.PutObjectOptions{ContentType: file.Header.Get("Content-Type")})
	// Expires: <-time.After(time.Hour * 24 * 30) <- I thought of keeping expiring objects but change of plans i guess

	if err != nil {
		log.Fatalln(err)
	}

	Tx, err := fh.db.BeginTx(context.Background(), &sql.TxOptions{})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed to start a PostgreSQL transaction!",
			"details": err.Error(),
		})
	}

	defer Tx.Rollback()

	query := `
		INSERT INTO gists (
		user_id,
		file_id,
		gist_title,
		short_url,
		is_public
		) VALUES ( $1, $2, $3, $4, $5 );
	`

	_, err = Tx.Exec(query,
		Gist.UserID,
		fileInfo.ETag,
		Gist.GistTitle,
		Gist.ShortUrl,
		Gist.IsPublic,
	)

	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed to insert gist into database!",
			"details": err.Error(),
		})
	}

	if err := Tx.Commit(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed to commit the PostgreSQL transaction!",
			"details": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]any{
		"key":      fileInfo.Key,
		"fileSize": fileInfo.Size,
		"fileTag":  fileInfo.ETag,
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
