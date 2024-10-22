package handlers

import (
	"bytes"
	"context"
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
	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
	"github.com/oklog/ulid/v2"
)

var MinioBucketName string = utils.GetEnv("MINIO_BUCKET_NAME", "codeflick")

type FileStorageHandler struct {
	gistDB  *db.GistDB
	replyDB *db.ReplyDB
	userDB  *db.UserDB
	minio   *db.MinioHandler
}

func NewFilesHandler(gistDB *db.GistDB, replyDB *db.ReplyDB, userDB *db.UserDB, minio *db.MinioHandler) *FileStorageHandler {
	return &FileStorageHandler{gistDB, replyDB, userDB, minio}
}

func (fh *FileStorageHandler) UploadGist(c echo.Context) error {
	var User models.User = c.Get("UserSessionDetails").(models.User)
	var Gist models.Gist
	Gist.UserID = User.ID
	currentTime := time.Now()
	Gist.CreatedAt = currentTime
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
		Gist.GistTitle = fmt.Sprintf("%s-%s", file.Filename, Gist.FileID)
	}

	Gist.ShortUrl = c.FormValue("custom_url")
	if Gist.ShortUrl == "" {
		Gist.ShortUrl = fmt.Sprintf("%s-%s", file.Filename, Gist.FileID)
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
	go func() {
		if err := src.Close(); err != nil {
			slog.AnyValue(fmt.Errorf("Error while closing the file: %w", err))
		}
	}()

	fileName := fmt.Sprintf("%s/%s", Gist.UserID, Gist.FileID)

	fileInfo, err := fh.minio.PutObject(context.Background(),
		fileName,
		src,
		file.Size,
		minio.PutObjectOptions{ContentType: file.Header.Get("Content-Type"), UserMetadata: map[string]string{"fileId": Gist.FileID}})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed to Upload File to minio!",
			"details": minio.ToErrorResponse(err),
		})
	}

	err = fh.gistDB.InsertGist(context.Background(), Gist)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			if err := fh.minio.RemoveObject(context.Background(), fileName, minio.RemoveObjectOptions{}); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
					"success": false,
					"message": "Failed to remove the file from minio!",
					"details": minio.ToErrorResponse(err),
				})
			}

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

	return c.JSON(http.StatusCreated, map[string]any{
		"success":    true,
		"message":    "Gist uploaded successfully!",
		"key":        fileInfo.Key,
		"fileSize":   fileInfo.Size,
		"file_id":    Gist.FileID,
		"short_url":  Gist.ShortUrl,
		"title":      Gist.GistTitle,
		"created_at": Gist.CreatedAt,
	})
}

func (fh *FileStorageHandler) ListBuckets(c echo.Context) error {
	var User models.User = c.Get("UserSessionDetails").(models.User)
	if !User.IsAdmin {
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

func (fh *FileStorageHandler) ListGists(c echo.Context) error {
	var User models.User = c.Get("UserSessionDetails").(models.User)

	gists, err := fh.gistDB.GetGistsByUserID(context.Background(), User.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed to get gists from database!",
			"details": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": "Fetch all gists successfully!",
		"data":    gists,
	})
}

func (fh *FileStorageHandler) GetGist(c echo.Context) error {
	GistUrl := c.Param("id")
	GistID := c.QueryParam("gid")
	if GistUrl == "" {
		return echo.NewHTTPError(http.StatusNotFound, map[string]any{
			"success": false,
			"message": "Gist ID is required! Not Found!",
		})
	}

	Gist, err := fh.gistDB.GetGistByShortURL(context.Background(), GistUrl)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, map[string]any{
			"success": false,
			"message": "Gist not found!",
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

	gistData, err := fh.minio.GetObject(context.Background(), fileName, minio.GetObjectOptions{})
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

	c.Response().Header().Set("Content-Type", "application/octet-stream")
	c.Response().Header().Set("X-Metadata-Length", strconv.Itoa(len(metadataJSON)))
	c.Response().Header().Set("X-Metadata", string(metadataJSON))
	c.Response().Header().Set("Content-Length", strconv.FormatInt(gistStat.Size, 10))

	var buffer bytes.Buffer

	if _, err := io.Copy(&buffer, gistData); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed to write file content",
			"details": err.Error(),
		})
	}

	return c.Blob(http.StatusOK, "application/octet-stream", buffer.Bytes())
}

func (fh *FileStorageHandler) UpdateGist(c echo.Context) error {
	currentTime := time.Now()
	var User models.User = c.Get("UserSessionDetails").(models.User)
	GistID := c.Param("id")

	if GistID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"success": false,
			"message": "Gist ID is required!",
		})
	}

	Gist, err := fh.gistDB.GetGistByID(context.Background(), GistID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, map[string]any{
			"success": false,
			"message": "Gist not found!",
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

	Gist.GistTitle = GistTitle
	Gist.ShortUrl = ShortUrl
	Gist.IsPublic = IsPublic
	Gist.UpdatedAt = currentTime

	returnGistId, err := fh.gistDB.UpdateGist(context.Background(), Gist)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return echo.NewHTTPError(http.StatusConflict, map[string]any{
				"success": false,
				"message": "Gist with the short/custom url already exists!",
				"details": err.Error(),
			})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed to update gist in database!",
			"details": err.Error(),
		})
	}

	file, err := c.FormFile("file")
	if err != nil {
		if strings.Contains(err.Error(), "no such file") {
			return c.JSON(http.StatusOK, map[string]any{
				"success": true,
				"message": "Gist Updated successfully!",
				"fileId":  returnGistId,
			})
		}
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"success": false,
			"message": "Error while reading the file!",
			"details": err.Error(),
		})
	}

	if !strings.Contains(file.Header.Get("Content-Type"), "text") {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"success": false,
			"message": "Only text files are allowed!",
		})
	}
	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Error While opening the file!",
			"details": err.Error(),
		})
	}
	go func() {
		if err := src.Close(); err != nil {
			slog.AnyValue(fmt.Errorf("Error while closing the file: %w", err))
		}
	}()

	fileName := fmt.Sprintf("%s/%s", Gist.UserID, Gist.FileID)

	_, err = fh.minio.PutObject(context.Background(),
		fileName,
		src,
		file.Size,
		minio.PutObjectOptions{ContentType: file.Header.Get("Content-Type"), UserMetadata: map[string]string{"fileId": Gist.FileID}})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed to Update File to minio!",
			"details": minio.ToErrorResponse(err),
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": "Gist Updated successfully!",
		"fileId":  returnGistId,
	})
}

func (fh *FileStorageHandler) DeleteGist(c echo.Context) error {
	var User models.User = c.Get("UserSessionDetails").(models.User)
	GistId := c.Param("id")

	if GistId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"success": false,
			"message": "Gist ID is required!",
		})
	}

	returnGistId, err := fh.gistDB.DeleteGist(context.Background(), GistId, User.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed to delete gist from database!",
			"details": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": "Gist deleted successfully!",
		"fileId":  returnGistId,
	})
}

func (fh *FileStorageHandler) GetGistReplies(c echo.Context) error {
	GistID := c.Param("id")

	slog.Debug("Fetching Gist Replies")

	if GistID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"success": false,
			"message": "Gist ID is required!",
		})
	}

	Gist, err := fh.gistDB.GetGistByID(context.Background(), GistID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, map[string]any{
			"success": false,
			"message": "Gist not found!",
			"details": err.Error(),
		})
	}

	if Gist.IsDeleted {
		return echo.NewHTTPError(http.StatusNotFound, map[string]any{
			"success": false,
			"message": "Gist not found! Deleted by the user!",
		})
	}

	replies, err := fh.replyDB.GetRepliesByGistID(context.Background(), GistID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed to get replies from database!",
			"details": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": "Replies fetched successfully!",
		"data":    replies,
	})
}

func (fh *FileStorageHandler) InsertGistReply(c echo.Context) error {
	currentTime := time.Now().UTC()
	GistID := c.Param("id")
	var User models.User = c.Get("UserSessionDetails").(models.User)
	var reply models.Reply = models.Reply{ID: ulid.Make().String(), UserID: User.ID, GistID: GistID, CreatedAt: currentTime, UpdatedAt: currentTime}

	if GistID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"success": false,
			"message": "Gist ID is required!",
		})
	}
	if err := c.Bind(&reply); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"success": false,
			"message": "Failed to decode request body!",
			"details": err.Error(),
		})
	}

	_, err := fh.gistDB.GetGistByID(context.Background(), GistID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, map[string]any{
			"success": false,
			"message": "Gist not found!",
			"details": err.Error(),
		})
	}

	err = fh.replyDB.InsertReply(context.Background(), reply)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed to add the reply!",
			"details": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": "Replied successfully!",
		"data":    reply,
	})
}

func (fh *FileStorageHandler) UpdateGistReply(c echo.Context) error {
	gistID := c.Param("id")
	replyID := c.Param("reply_id")
	currentTime := time.Now().UTC()
	var User models.User = c.Get("UserSessionDetails").(models.User)
	var reply models.Reply = models.Reply{ID: replyID, GistID: gistID, UserID: User.ID, UpdatedAt: currentTime}

	if gistID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"success": false,
			"message": "Gist ID is required!",
		})
	}

	if replyID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"success": false,
			"message": "Reply ID is required!",
		})
	}

	if err := c.Bind(&reply); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"success": false,
			"message": "Failed to decode request body!",
			"details": err.Error(),
		})
	}

	existingReply, err := fh.replyDB.GetReplyByID(context.Background(), replyID, gistID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, map[string]any{
			"success": false,
			"message": "Reply not found!",
			"details": err.Error(),
		})
	}

	if existingReply.UserID != User.ID {
		return echo.NewHTTPError(http.StatusUnauthorized, map[string]any{
			"success": false,
			"message": "You are not allowed to update others reply!",
		})
	}

	if existingReply.GistID != gistID {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"success": false,
			"message": "Reply does not belong to the Gist!",
		})
	}

	if existingReply.IsDeleted {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"success": false,
			"message": "Reply is Deleted! Can't Update it!",
		})
	}

	err = fh.replyDB.UpdateReply(context.Background(), reply)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed to update the reply!",
			"details": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": "Updated Reply successfully!",
		"data":    reply,
	})
}

func (fh *FileStorageHandler) DeleteGistReply(c echo.Context) error {
	gistID := c.Param("id")
	replyID := c.Param("reply_id")
	var User models.User = c.Get("UserSessionDetails").(models.User)

	if gistID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"success": false,
			"message": "Gist ID is required!",
		})
	}

	if replyID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"success": false,
			"message": "Reply ID is required!",
		})
	}

	reply, err := fh.replyDB.GetReplyByID(context.Background(), replyID, gistID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, map[string]any{
			"success": false,
			"message": "Reply not found!",
			"details": err.Error(),
		})
	}

	if reply.UserID != User.ID {
		return echo.NewHTTPError(http.StatusUnauthorized, map[string]any{
			"success": false,
			"message": "You are not allowed to delete others reply!",
		})
	}

	if reply.GistID != gistID {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"success": false,
			"message": "Reply does not belong to the Gist!",
		})
	}

	if reply.IsDeleted {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"success": false,
			"message": "Reply was already deleted!",
		})
	}

	err = fh.replyDB.DeleteReply(context.Background(), replyID, User.ID, gistID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed to delete the reply!",
			"details": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": "Deleted Reply successfully!",
	})
}

func (fh *FileStorageHandler) ListAllFiles(c echo.Context) error {
	return nil
}
