package routes

import (
	"net/http"
	"time"

	"github.com/HarshPatel5940/CodeFlick/internal/handlers"
	m "github.com/HarshPatel5940/CodeFlick/internal/middlewares"
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
		"uptime": time.Since(StartTime).String(),
	})
}

func SetupAPIRoutes(e *echo.Group,
	gh *handlers.GistHandler,
	ah *handlers.AuthHandler,
) {
	e.GET("", root)

	e.GET("/auth/:provider/login", ah.GoogleOauthLogin)
	e.GET("/auth/:provider/callback", ah.GoogleOauthCallback)
	e.GET("/auth/session", ah.GetSessionDetails, m.SessionMiddleware(m.Config{
		RequiredAccess: m.AccessLevelUser,
	}))
	e.POST("/auth/logout", ah.Logout)

	gistsRoutes := e.Group("/gists", m.SessionMiddleware(m.Config{
		RequiredAccess: m.AccessLevelUser,
	}))
	gistsRoutes.GET("", gh.GetAllGists)
	gistsRoutes.POST("/new", gh.UploadGist)
	gistsRoutes.GET("/:id", gh.GetGist)
	gistsRoutes.PUT("/:id", gh.UpdateGist)
	gistsRoutes.DELETE("/:id", gh.DeleteGist)

	gistsRoutes.GET("/:id/reply", gh.GetGistReplies)
	gistsRoutes.POST("/:id/reply", gh.InsertGistReply)
	gistsRoutes.PUT("/:id/reply/:reply_id", gh.UpdateGistReply)
	gistsRoutes.DELETE("/:id/reply/:reply_id", gh.DeleteGistReply)

	adminRoutes := e.Group("/admin", m.SessionMiddleware(m.Config{
		RequiredAccess: m.AccessLevelAdmin,
	}))
	adminRoutes.GET("/buckets", gh.ListBuckets)
	adminRoutes.GET("/buckets/:bucket", gh.ListAllFiles)
	adminRoutes.GET("/users", ah.GetAllUsers)
	adminRoutes.PUT("/gists/:id/reply/:reply_id", gh.UpdateGist)
	adminRoutes.DELETE("/gists/:id/reply/:reply_id", gh.DeleteGist)
}

func SetupPageRoutes(e *echo.Group) {
	e.GET("", func(c echo.Context) error {
		return pages.Render(c, http.StatusOK, pages.Home())
	})
}
