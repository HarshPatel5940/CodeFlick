package middlewares

import (
	"fmt"
	"net/http"

	"github.com/HarshPatel5940/CodeFlick/internal/utils"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	m "github.com/labstack/echo/v4/middleware"
)

func SetupMiddlewares(app *echo.Echo) {
	app.Pre(m.RemoveTrailingSlash())
	app.Use(m.Logger())
	app.Use(m.Decompress())
	app.Use(m.CORSWithConfig(m.CORSConfig{
		AllowOrigins: []string{
			fmt.Sprintf("localhost:%s", utils.GetEnv("PORT")),
			utils.GetEnv("CLIENT_URL"),
		},
		AllowCredentials: true,
		ExposeHeaders: []string{"gist-metadata", "Gist-Metadata"}, // Try both cases
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	app.Use(m.Secure())
	app.Use(session.Middleware(sessions.NewCookieStore([]byte(utils.GetEnv("GORILLA_SESSIONS_KEY")))))
	app.Use(m.RateLimiter(m.NewRateLimiterMemoryStore(20)))
	app.Use(m.Recover())
}
