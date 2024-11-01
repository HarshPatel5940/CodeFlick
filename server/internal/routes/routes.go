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
	GistStorageHandler *handlers.GistStorageHandler,
	AuthHandler *handlers.AuthHandler,
) {
	e.GET("", root)

	e.GET("/auth/:provider/login", AuthHandler.GoogleOauthLogin)
	e.GET("/auth/:provider/callback", AuthHandler.GoogleOauthCallback)
	e.GET("/auth/session", AuthHandler.GetSessionDetails, middlewares.SessionMiddleware(middlewares.Config{
		RequiredAccess: middlewares.AccessLevelUser,
	}))

	gistsRoutes := e.Group("/gists", middlewares.SessionMiddleware(middlewares.Config{
		RequiredAccess: middlewares.AccessLevelUser,
	}))
	gistsRoutes.GET("", GistStorageHandler.GetAllGists)
	gistsRoutes.POST("/new", GistStorageHandler.UploadGist)
	gistsRoutes.GET("/:id", GistStorageHandler.GetGist)
	gistsRoutes.PUT("/:id", GistStorageHandler.UpdateGist)
	gistsRoutes.DELETE("/:id", GistStorageHandler.DeleteGist)

	gistsRoutes.GET("/:id/reply", GistStorageHandler.GetGistReplies)
	gistsRoutes.POST("/:id/reply", GistStorageHandler.InsertGistReply)
	gistsRoutes.PUT("/:id/reply/:reply_id", GistStorageHandler.UpdateGistReply)
	gistsRoutes.DELETE("/:id/reply/:reply_id", GistStorageHandler.DeleteGistReply)

	adminRoutes := e.Group("/admin", middlewares.SessionMiddleware(middlewares.Config{
		RequiredAccess: middlewares.AccessLevelAdmin,
	}))
	adminRoutes.GET("/buckets", GistStorageHandler.ListBuckets)
	adminRoutes.GET("/buckets/:bucket", GistStorageHandler.ListAllFiles)
	adminRoutes.GET("/users", AuthHandler.GetAllUsers)
	// TODO: Did not test this, so do it when u get time.
	adminRoutes.PUT("/gists/:id/reply/:reply_id", GistStorageHandler.UpdateGist)
	adminRoutes.DELETE("/gists/:id/reply/:reply_id", GistStorageHandler.DeleteGist)
}

func SetupPageRoutes(e *echo.Group) {
	e.GET("", func(c echo.Context) error {
		return pages.Render(c, http.StatusOK, pages.Home())
	})
}
