package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
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
	if err != nil || sess.Values["user_id"] == nil || sess.Values["user_id"] == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, map[string]any{
			"success": false,
			"message": "Failed to get session!",
		})
	}

	var Gist models.Gist
	Gist.UserID = fmt.Sprint(sess.Values["user_id"])
	currentTime := time.Now().Unix()
	Gist.FileID = ulid.Make().String()

	file, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"success": false,
			"message": "No File Provided for the gist!",
			"details": err.Error(),
		})
	}

	if !strings.Contains(file.Header.Get("Content-Type"), "text") {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"success": false,
			"message": "Only text files are allowed!",
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

	fileName := fmt.Sprintf("%s/%s", Gist.UserID, Gist.FileID)

	fileInfo, err := fh.minio.PutObject(context.Background(),
		MinioBucketName,
		fileName,
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
	sess, err := session.Get("session", c)

	if err != nil || sess.Values["user_id"] == nil || sess.Values["user_id"] == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, map[string]any{
			"success": false,
			"message": "Failed to get session!",
		})
	}

	IsAdmin, err := strconv.ParseBool(c.FormValue("is_public"))
	if err != nil {
		IsAdmin = false
	}

	if !IsAdmin {
		return echo.NewHTTPError(http.StatusUnauthorized, map[string]any{
			"success": false,
			"message": "You are not authorized to view this page!",
		})
	}

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
	var Gist models.Gist
	GistUrl := c.Param("id")
	GistID := c.QueryParam("gid")
	if GistUrl == "" {
		return echo.NewHTTPError(http.StatusNotFound, map[string]any{
			"success": false,
			"message": "Gist ID is required! Not Found!",
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

	orgGistRow := fh.db.QueryRowContext(context.Background(), db.GetGistByShortURL, GistUrl)
	err = orgGistRow.Scan(&Gist.FileID, &Gist.UserID, &Gist.ForkedFrom, &Gist.GistTitle, &Gist.ShortUrl,
		&Gist.ViewCount, &Gist.IsPublic, &Gist.IsDeleted, &Gist.CreatedAt, &Gist.UpdatedAt)

	if err := Tx.Commit(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed to commit the PostgreSQL transaction!",
			"details": err.Error(),
		})
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, map[string]any{
				"success": false,
				"message": "Gist not found!",
				"details": err.Error(),
			})
		}
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed to get gist from database!",
			"details": err.Error(),
		})
	}

	if Gist.IsDeleted {
		return echo.NewHTTPError(http.StatusNotFound, map[string]any{
			"success": false,
			"message": "Gist not found!",
		})
	}

	if !Gist.IsPublic {
		if Gist.FileID != GistID {
			return echo.NewHTTPError(http.StatusUnauthorized, map[string]any{
				"success": false,
				"message": "You are not authorized to get this gist!",
			})
		}
	}
	fileName := fmt.Sprintf("%s/%s", Gist.UserID, Gist.FileID)

	gistData, err := fh.minio.GetObject(context.Background(), MinioBucketName, fileName, minio.GetObjectOptions{})
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

	gistStat, err := gistData.Stat()
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

	metadataJSON, err := json.Marshal(map[string]interface{}{
		"UserID":     Gist.UserID,
		"ForkedFrom": Gist.ForkedFrom,
		"GistTitle":  Gist.GistTitle,
		"ShortUrl":   Gist.ShortUrl,
		"ViewCount":  Gist.ViewCount,
		"IsPublic":   Gist.IsPublic,
		"CreatedAt":  Gist.CreatedAt,
		"UpdatedAt":  Gist.UpdatedAt,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed to marshal metadata",
			"details": err.Error(),
		})
	}

	// Set response headers
	c.Response().Header().Set("Content-Type", "application/octet-stream")
	c.Response().Header().Set("X-Metadata-Length", strconv.Itoa(len(metadataJSON)))
	c.Response().Header().Set("X-Metadata", string(metadataJSON))
	c.Response().Header().Set("Content-Length", strconv.FormatInt(gistStat.Size, 10))

	// Create a buffer to hold the full response
	var buffer bytes.Buffer

	if _, err := io.Copy(&buffer, gistData); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed to write file content",
			"details": err.Error(),
		})
	}

	// Return the response using c.Blob()
	return c.Blob(http.StatusOK, "application/octet-stream", buffer.Bytes())
}

func (fh FileStorageHandler) UpdateGist(c echo.Context) error {
	currentTime := time.Now()
	var Gist models.Gist
	var User models.User = c.Get("UserSessionDetails").(models.User)
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

	orgGistRow := Tx.QueryRowContext(context.Background(), db.GetGistByID, Gist.FileID)
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

	if Gist.UserID != User.ID {
		return echo.NewHTTPError(http.StatusForbidden, map[string]any{
			"success": false,
			"message": "You are not allowed to update this gist!",
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
	UserID := fmt.Sprint(sess.Values["user_id"])

	if err != nil || sess.Values["user_id"] == nil || sess.Values["user_id"] == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, map[string]any{
			"success": false,
			"message": "Failed to get session!",
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

	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, map[string]any{
				"success": false,
				"message": "Failed to update gist from database! Not Found",
				"details": err.Error(),
			})
		}

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
