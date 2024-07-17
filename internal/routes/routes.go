package routes

import (
	"net/http"
	"time"

	"github.com/HarshPatel5940/CodeFlick/internal/handlers"
	"github.com/labstack/echo/v4"
)

var StartTime time.Time

func root(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{
		"message": []string{
			"Hello, World! This is CodeFlick API",
			"Codeflick is an pastebin alternative to distribute your gists / Code quickly!",
		},
		"my-links": map[string]string{
			"me":      "https://github.com/HarshPatel5940",
			"twitter": "https://twitter.com/HarshPatel5940",
			"project": "https://github.com/HarshPatel5940/CodeFlick",
		},
	})
}

func SetupRoutes(e *echo.Group,
	FileStorageHandler *handlers.FileStorageHandler,
	AuthHandler *handlers.AuthHandler,
) {
	e.GET("", root)

	// Auth Routes
	e.GET("/auth/:provider/login", AuthHandler.GoogleOauthLogin)
	e.GET("/auth/:provider/callback", AuthHandler.GoogleOauthCallback)
	e.GET("/auth/session", AuthHandler.GetSessionDetails)

	// File Routes
	e.POST("/gists/new", FileStorageHandler.UploadFile)
	e.GET("/admin/buckets", FileStorageHandler.ListBuckets)

}
