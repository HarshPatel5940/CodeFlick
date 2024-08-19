package middlewares

import (
	"net/http"

	"github.com/HarshPatel5940/CodeFlick/internal/models"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func SessionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("session", c)

		if err != nil || sess.Values["user_id"] == nil || sess.Values["user_id"] == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, map[string]any{
				"success": false,
				"message": "Failed to get session!",
			})
		}

		var User models.User = models.User{
			ID:           sess.Values["user_id"].(string),
			Name:         sess.Values["name"].(string),
			Email:        sess.Values["email"].(string),
			AuthProvider: sess.Values["auth_provider"].(string),
			IsAdmin:      sess.Values["isAdmin"].(bool),
			IsPremium:    sess.Values["isPremium"].(bool),
			IsDeleted:    sess.Values["isDeleted"].(bool),
		}

		c.Set("UserSessionDetails", User)

		return next(c)
	}
}
