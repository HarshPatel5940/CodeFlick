package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/HarshPatel5940/CodeFlick/internal/db"
	"github.com/HarshPatel5940/CodeFlick/internal/models"
	"github.com/HarshPatel5940/CodeFlick/internal/utils"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
	"github.com/oklog/ulid/v2"
)

var MinioBucketName string = utils.GetEnv("MINIO_BUCKET_NAME", "codeflick")

type FileStorageHandler struct {
	db    *sqlx.DB
	minio *minio.Client
}

func NewFilesHandler(db *sqlx.DB, minio *minio.Client) *FileStorageHandler {
	return &FileStorageHandler{db, minio}
}

func (fh FileStorageHandler) UploadGist(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, map[string]any{
			"success": false,
			"message": "Failed to get session!",
			"details": err.Error(),
		})
	}

	var Gist models.Gist
	Gist.UserID = sess.Values["user_id"].(string)
	currentTime := time.Now().Unix()

	file, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"success": false,
			"message": "No File Provided for the gist!",
			"details": err.Error(),
		})
	}

	Gist.GistTitle = c.FormValue("gist_title")
	if Gist.GistTitle == "" {
		Gist.GistTitle = fmt.Sprintf("%s-%d", file.Filename, currentTime)
	}

	Gist.ShortUrl = c.FormValue("custom_url")
	if Gist.ShortUrl == "" {
		Gist.ShortUrl = fmt.Sprintf("%s-%d", file.Filename, currentTime)
	}

	Gist.IsPublic, err = strconv.ParseBool(c.FormValue("is_public"))
	if err != nil {
		Gist.IsPublic = false
	}

	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Error While opening the file!",
			"details": err.Error(),
		})
	}
	defer src.Close()

	fileNameParts := strings.SplitAfter(file.Filename, ".")

	location := fmt.Sprintf("%s/%d.%s",
		Gist.UserID, currentTime, fileNameParts[len(fileNameParts)-1],
	)

	Gist.FileID = ulid.Make().String()
	fileInfo, err := fh.minio.PutObject(context.Background(),
		MinioBucketName,
		location,
		src,
		file.Size,
		minio.PutObjectOptions{ContentType: file.Header.Get("Content-Type"), UserMetadata: map[string]string{"fileId": Gist.FileID}})
	// Expires: <-time.After(time.Hour * 24 * 30) <- I thought of keeping expiring objects but change of plans i guess

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed to Upload File to minio!",
			"details": minio.ToErrorResponse(err),
		})
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

	_, err = Tx.Exec(db.InsertGist, Gist.UserID, Gist.FileID, Gist.GistTitle, Gist.ShortUrl, Gist.IsPublic)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			fh.minio.RemoveObject(context.Background(), MinioBucketName, fileInfo.Key, minio.RemoveObjectOptions{})

			return echo.NewHTTPError(http.StatusConflict, map[string]any{
				"success": false,
				"message": "Gist with the same title or short url already exists!",
				"details": err.Error(),
			})
		}
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
		"success":  true,
		"message":  "Gist uploaded successfully!",
		"key":      fileInfo.Key,
		"fileSize": fileInfo.Size,
		"fileId":   Gist.FileID,
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

func (fh FileStorageHandler) GetGist(c echo.Context) error {
	// TODO: Complete this logic with public / private logic lateer
	return nil
}

func (fh FileStorageHandler) GetGistRaw(c echo.Context) error {
	// TODO: Complete this logic with public / private logic lateer
	return nil
}

func (fh FileStorageHandler) UpdateGist(c echo.Context) error {
	sess, err := session.Get("session", c)
	currentTime := time.Now()
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, map[string]any{
			"success": false,
			"message": "Failed to get session!",
			"details": err.Error(),
		})
	}

	var Gist models.Gist
	Gist.UserID = sess.Values["user_id"].(string)
	Gist.FileID = c.Param("id")

	if Gist.FileID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"success": false,
			"message": "Gist ID is required!",
		})
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

	orgGistRow := Tx.QueryRowContext(context.Background(), db.GetGistByID, Gist.FileID, Gist.UserID)
	err = orgGistRow.Scan(&Gist.FileID, &Gist.UserID, &Gist.ForkedFrom, &Gist.GistTitle, &Gist.ShortUrl,
		&Gist.ViewCount, &Gist.IsPublic, &Gist.IsDeleted, &Gist.CreatedAt, &Gist.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, map[string]any{
				"success": false,
				"message": "Gist not found!",
			})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed to get gist from database!",
			"details": err.Error(),
		})
	}

	GistTitle := c.FormValue("gist_title")
	if GistTitle == "" {
		GistTitle = Gist.GistTitle
	}

	ShortUrl := c.FormValue("custom_url")
	if ShortUrl == "" {
		ShortUrl = Gist.ShortUrl
	}

	IsPublic, err := strconv.ParseBool(c.FormValue("is_public"))
	if err != nil {
		IsPublic = Gist.IsPublic
	}

	slog.Info(Gist.FileID, Gist.UserID, GistTitle, ShortUrl, IsPublic)

	var returnGistId string
	row := Tx.QueryRowContext(context.Background(), db.UpdateGist, Gist.FileID, Gist.UserID, GistTitle, ShortUrl, IsPublic, currentTime)
	err = row.Scan(&returnGistId)

	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, map[string]any{
				"success": false,
				"message": "Failed to update gist from database! Not Found",
				"details": err.Error(),
			})
		}

		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return echo.NewHTTPError(http.StatusConflict, map[string]any{
				"success": false,
				"message": "Gist with the short/custom url already exists!",
				"details": err.Error(),
			})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed to update gist from database!",
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

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": "Gist Updated successfully!",
		"fileId":  returnGistId,
	})
}

func (fh FileStorageHandler) DeleteGist(c echo.Context) error {
	sess, err := session.Get("session", c)
	UserID := sess.Values["user_id"].(string)

	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, map[string]any{
			"success": false,
			"message": "Failed to get session!",
			"details": err.Error(),
		})
	}
	GistId := c.Param("id")
	var returnGistId string

	if GistId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"success": false,
			"message": "Gist ID is required!",
		})
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

	row := Tx.QueryRowContext(context.Background(), db.DeleteGist, GistId, UserID, time.Now())
	err = row.Scan(&returnGistId)

	if err != nil || returnGistId == "" {
		// TODO: HANDLE logic for not found error
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed to delete gist from database!",
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

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": "Gist deleted successfully!",
		"fileId":  returnGistId,
	})
}
