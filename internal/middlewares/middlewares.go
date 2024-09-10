package middlewares

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupMiddlewares(app *echo.Echo) {
	app.Pre(middleware.RemoveTrailingSlash())
	app.Use(middleware.Logger())
	app.Use(middleware.Decompress())
	app.Use(middleware.CORS())
	app.Use(middleware.Secure())
	app.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))
	app.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(4)))
	app.Use(middleware.Recover())
}
