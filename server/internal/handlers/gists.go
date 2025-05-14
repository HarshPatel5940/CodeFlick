package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/HarshPatel5940/CodeFlick/internal/db"
	"github.com/HarshPatel5940/CodeFlick/internal/models"
	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
	"github.com/oklog/ulid/v2"
)

type GistHandler struct {
	gistDB  *db.GistDB
	replyDB *db.ReplyDB
	userDB  *db.UserDB
	minio   *db.MinioHandler
}

func NewGistHandler(gistDB *db.GistDB, replyDB *db.ReplyDB, userDB *db.UserDB, minio *db.MinioHandler) *GistHandler {
	return &GistHandler{gistDB, replyDB, userDB, minio}
}

func (gh *GistHandler) UploadGist(c echo.Context) error {
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

	fileNameReq := c.FormValue("file_name")
	if fileNameReq == "" {
		fileNameReq = file.Filename
	}

	Gist.FileName = fmt.Sprintf("%s/%s-%s", Gist.UserID, Gist.FileID, fileNameReq)
	if len(Gist.FileName) > 100 {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"success": false,
			"message": "File name is too long! Max length is 100 characters.",
		})
	}

	Placeholder := fmt.Sprintf("%s-%s", Gist.FileName, Gist.FileID)

	if !strings.Contains(file.Header.Get("Content-Type"), "text") {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"success": false,
			"message": "Only text files are allowed!",
		})
	}

	Gist.GistTitle = c.FormValue("gist_title")
	if Gist.GistTitle == "" {
		Gist.GistTitle = Placeholder
	}

	if len(Gist.GistTitle) < 5 {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"success": false,
			"message": "Gist Title must be atleast 5 characters long!",
		})
	}

	Gist.ShortUrl = c.FormValue("custom_url")
	if Gist.ShortUrl == "" {
		Gist.ShortUrl = Placeholder
	}

	if len(Gist.ShortUrl) < 5 {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"success": false,
			"message": "Short URL must be atleast 5 characters long!",
		})
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
			slog.AnyValue(fmt.Errorf("error while closing the file: %w", err))
		}
	}()

	_, err = gh.minio.PutObject(context.Background(),
		Gist.FileName,
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

	_, err = gh.gistDB.InsertGist(context.Background(), Gist)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			if err := gh.minio.RemoveObject(context.Background(), Gist.FileName, minio.RemoveObjectOptions{}); err != nil {
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
		"file_id":    Gist.FileID,
		"short_url":  Gist.ShortUrl,
		"title":      Gist.GistTitle,
		"created_at": Gist.CreatedAt,
	})
}

func (gh *GistHandler) GetAllGists(c echo.Context) error {
	var User models.User = c.Get("UserSessionDetails").(models.User)
	fetchPublic := c.QueryParam("fetchPublic")

	gists, err := gh.gistDB.GetGistsByUserID(context.Background(), User.ID, fetchPublic)
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

func (gh *GistHandler) GetGist(c echo.Context) error {
	GistUrl := c.Param("id")
	GistID := c.QueryParam("gid")
	if GistUrl == "" {
		return echo.NewHTTPError(http.StatusNotFound, map[string]any{
			"success": false,
			"message": "Gist URL is required",
		})
	}

	Gist, err := gh.gistDB.GetGistByShortURL(context.Background(), GistUrl)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, map[string]any{
			"success": false,
			"message": "Gist not found",
		})
	}

	if Gist.IsDeleted {
		return echo.NewHTTPError(http.StatusNotFound, map[string]any{
			"success": false,
			"message": "Gist has been deleted",
		})
	}

	if !Gist.IsPublic && Gist.FileID != GistID {
		return echo.NewHTTPError(http.StatusUnauthorized, map[string]any{
			"success": false,
			"message": "Unauthorized access to private gist",
		})
	}

	gistData, err := gh.minio.GetObject(context.Background(), Gist.FileName, minio.GetObjectOptions{})
	if err != nil {
		c.Logger().Error(err, gistData
		)
		return handleMinioError(err)
	}
	defer func() {
		if err := gistData.Close(); err != nil {
			slog.AnyValue(fmt.Errorf("error while closing the file: %w", err))
		}
	}()

	metadata := map[string]interface{}{
		"userID":     Gist.UserID,
		"fileID":     Gist.FileID,
		"fileName":   Gist.FileName,
		"gistTitle":  Gist.GistTitle,
		"shortUrl":   Gist.ShortUrl,
		"viewCount":  Gist.ViewCount,
		"isPublic":   Gist.IsPublic,
		"createdAt":  Gist.CreatedAt,
		"updatedAt":  Gist.UpdatedAt,
		"forkedFrom": Gist.ForkedFrom,
	}

	buffer := &bytes.Buffer{}
	if _, err := io.Copy(buffer, gistData); err != nil {
		c.Logger().Error(err, Gist)
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed to read file content",
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"success":  true,
		"message":  "Gist fetched successfully",
		"content":  buffer.String(),
		"metadata": metadata,
	})
}

func (gh *GistHandler) UpdateGist(c echo.Context) error {
	currentTime := time.Now()
	var User models.User = c.Get("UserSessionDetails").(models.User)
	GistID := c.Param("id")

	if GistID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"success": false,
			"message": "Gist ID is required!",
		})
	}

	Gist, err := gh.gistDB.GetGistByID(context.Background(), GistID)
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

	returnGistId, err := gh.gistDB.UpdateGist(context.Background(), Gist)
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
			slog.AnyValue(fmt.Errorf("error while closing the file: %w", err))
		}
	}()

	originalFileName := Gist.FileName
	newFileName := fmt.Sprintf("%s/%s-%s", Gist.UserID, Gist.FileID, file.Filename)

	Gist.FileName = newFileName
	_, err = gh.gistDB.UpdateGistFileName(context.Background(), Gist)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed to update gist filename in database!",
			"details": err.Error(),
		})
	}

	_, err = gh.minio.PutObject(context.Background(),
		newFileName,
		src,
		file.Size,
		minio.PutObjectOptions{ContentType: file.Header.Get("Content-Type"), UserMetadata: map[string]string{"fileId": Gist.FileID}})
	if err != nil {
		Gist.FileName = originalFileName
		_, revertErr := gh.gistDB.UpdateGistFileName(context.Background(), Gist)
		if revertErr != nil {
			slog.AnyValue(fmt.Errorf("failed to revert filename change in database: %w", revertErr))
		}

		return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed to Update File to minio!",
			"details": minio.ToErrorResponse(err),
		})
	}

	if originalFileName != newFileName {
		if errRemove := gh.minio.RemoveObject(context.Background(), originalFileName, minio.RemoveObjectOptions{}); errRemove != nil {
			slog.AnyValue(fmt.Errorf("failed to remove old file from minio: %w", errRemove))
		}
	}

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": "Gist Updated successfully!",
		"fileId":  returnGistId,
	})
}

func (gh *GistHandler) DeleteGist(c echo.Context) error {
	var User models.User = c.Get("UserSessionDetails").(models.User)
	GistId := c.Param("id")

	if GistId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"success": false,
			"message": "Gist ID is required!",
		})
	}

	returnGistId, err := gh.gistDB.DeleteGist(context.Background(), GistId, User.ID)
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

func (gh *GistHandler) GetGistReplies(c echo.Context) error {
	GistID := c.Param("id")

	slog.Debug("Fetching Gist Replies")

	if GistID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"success": false,
			"message": "Gist ID is required!",
		})
	}

	Gist, err := gh.gistDB.GetGistByID(context.Background(), GistID)
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

	replies, err := gh.replyDB.GetRepliesByGistID(context.Background(), GistID)
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

func (gh *GistHandler) InsertGistReply(c echo.Context) error {
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

	_, err := gh.gistDB.GetGistByID(context.Background(), GistID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, map[string]any{
			"success": false,
			"message": "Gist not found!",
			"details": err.Error(),
		})
	}

	err = gh.replyDB.InsertReply(context.Background(), reply)
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

func (gh *GistHandler) UpdateGistReply(c echo.Context) error {
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

	existingReply, err := gh.replyDB.GetReplyByID(context.Background(), replyID, gistID)
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

	err = gh.replyDB.UpdateReply(context.Background(), reply)
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

func (gh *GistHandler) DeleteGistReply(c echo.Context) error {
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

	reply, err := gh.replyDB.GetReplyByID(context.Background(), replyID, gistID)
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

	err = gh.replyDB.DeleteReply(context.Background(), replyID, User.ID, gistID)
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

func (gh *GistHandler) ListAllFiles(c echo.Context) error {
	gistID := c.Param("id")
	if gistID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"success": false,
			"message": "Gist ID is required!",
		})
	}

	name := gh.minio.GetBucketName()

	var fileObjects []minio.ObjectInfo
	for object := range gh.minio.ListObjects(context.Background(), name, true) {
		if object.Err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
				"success": false,
				"message": "Error listing objects",
				"details": object.Err.Error(),
			})
		}
		fileObjects = append(fileObjects, object)
	}

	if len(fileObjects) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, map[string]any{
			"success": false,
			"message": "No files found!",
		})
	}

	files := fileObjects

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": "Listed files successfully!",
		"data":    files,
	})
}

func (gh *GistHandler) ListBuckets(c echo.Context) error {
	var User models.User = c.Get("UserSessionDetails").(models.User)
	if !User.IsAdmin {
		return echo.NewHTTPError(http.StatusUnauthorized, map[string]any{
			"success": false,
			"message": "You are not authorized to view this page!",
		})
	}

	list, err := gh.minio.ListBuckets(context.Background())
	if err != nil {
		return handleMinioError(err)
	}

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": list,
	})
}

func handleMinioError(err error) error {
	if resp := minio.ToErrorResponse(err); resp.Code != "" {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": resp.Message,
			"code":    resp.Code,
		})
	}
	return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
		"success": false,
		"message": "Storage error occurred",
	})
}
