package middlewares

import (
	"fmt"

	"github.com/HarshPatel5940/CodeFlick/internal/utils"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupMiddlewares(app *echo.Echo) {
	app.Pre(middleware.RemoveTrailingSlash())
	app.Use(middleware.Logger())
	app.Use(middleware.Decompress())
	app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			fmt.Sprintf("localhost:%s", utils.GetEnv("PORT")),
			utils.GetEnv("CLIENT_URL"),
		},
	}))
	app.Use(middleware.Secure())
	app.Use(session.Middleware(sessions.NewCookieStore([]byte(utils.GetEnv("GORILLA_SESSIONS_KEY")))))
	app.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(10)))
	app.Use(middleware.Recover())
}
