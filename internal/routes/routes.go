package routes

import (
	"net/http"
	"time"

	"github.com/HarshPatel5940/CodeFlick/internal/handlers"
	"github.com/HarshPatel5940/CodeFlick/internal/middlewares"
	"github.com/HarshPatel5940/CodeFlick/internal/pages"
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

func SetupAPIRoutes(e *echo.Group,
	FileStorageHandler *handlers.FileStorageHandler,
	AuthHandler *handlers.AuthHandler,
) {
	e.GET("", root)

	// * ==========================
	e.GET("/auth/:provider/login", AuthHandler.GoogleOauthLogin, middlewares.SessionMiddleware(middlewares.Config{}))
	e.GET("/auth/:provider/callback", AuthHandler.GoogleOauthCallback, middlewares.SessionMiddleware(middlewares.Config{}))
	e.GET("/auth/session", AuthHandler.GetSessionDetails, middlewares.SessionMiddleware(middlewares.Config{
		RequiredAccess: middlewares.AccessLevelUser,
	}))

	// * ==========================
	gistsRoutes := e.Group("/gists", middlewares.SessionMiddleware(middlewares.Config{
		RequiredAccess: middlewares.AccessLevelUser,
	}))
	gistsRoutes.GET("", FileStorageHandler.ListGists)
	gistsRoutes.POST("/new", FileStorageHandler.UploadGist)
	gistsRoutes.GET("/:id", FileStorageHandler.GetGist)
	gistsRoutes.PUT("/:id", FileStorageHandler.UpdateGist)
	gistsRoutes.DELETE("/:id", FileStorageHandler.DeleteGist)

	gistsRoutes.GET("/:id/reply", FileStorageHandler.GetGistReplies)
	gistsRoutes.POST("/:id/reply", FileStorageHandler.InsertGistReply)
	gistsRoutes.PUT("/:id/reply/:reply_id", FileStorageHandler.UpdateGistReply)
	gistsRoutes.DELETE("/:id/reply/:reply_id", FileStorageHandler.DeleteGistReply)

	// * ========================== Admin Routes ==================
	adminRoutes := e.Group("/admin", middlewares.SessionMiddleware(middlewares.Config{
		RequiredAccess: middlewares.AccessLevelAdmin,
	}))
	adminRoutes.GET("/buckets", FileStorageHandler.ListBuckets)
	adminRoutes.GET("/buckets/:bucket", FileStorageHandler.ListAllFiles)
	// TODO: adminRoutes.GET("/users", FileStorageHandler.GetAllUsers)
	// TODO: delete and update routes for users and gists
}

func SetupPageRoutes(e *echo.Group) {
	e.GET("", func(c echo.Context) error {

		return pages.Render(c, http.StatusOK, pages.Home())
	})
}
