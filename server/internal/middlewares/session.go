package middlewares

import (
	"log/slog"
	"net/http"

	"github.com/HarshPatel5940/CodeFlick/internal/models"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type AccessLevel string

const (
	AccessLevelAll   AccessLevel = "all"
	AccessLevelUser  AccessLevel = "user"
	AccessLevelAdmin AccessLevel = "admin"
)

type Config struct {
	RequiredAccess AccessLevel
}

func GetSessionValue[T any](s *sessions.Session, key string, defaultValue T) T {
	if value, ok := s.Values[key]; ok {
		if typedValue, ok := value.(T); ok {
			return typedValue
		}
	}
	return defaultValue
}

func SessionMiddleware(config Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess, err := session.Get("_gothic_session", c)

			if config.RequiredAccess == "" {
				config.RequiredAccess = AccessLevelAll
			}

			if err != nil {
				if config.RequiredAccess != AccessLevelAll {
					return echo.NewHTTPError(http.StatusUnauthorized, map[string]any{
						"success": false,
						"message": "Failed to get session!",
					})
				}
			}

			userID := GetSessionValue(sess, "user_id", "")
			if userID == "" && config.RequiredAccess != AccessLevelAll {
				return echo.NewHTTPError(http.StatusUnauthorized, map[string]any{
					"success": false,
					"message": "User not authenticated",
				})
			}

			User := models.User{
				ID:           userID,
				Name:         GetSessionValue(sess, "name", ""),
				Email:        GetSessionValue(sess, "email", ""),
				AuthProvider: GetSessionValue(sess, "auth_provider", ""),
				IsAdmin:      GetSessionValue(sess, "isAdmin", false),
				IsPremium:    GetSessionValue(sess, "isPremium", false),
				IsDeleted:    GetSessionValue(sess, "isDeleted", false),
			}

			if config.RequiredAccess == AccessLevelAdmin && !User.IsAdmin {
				return echo.NewHTTPError(http.StatusUnauthorized, map[string]any{
					"success": false,
					"message": "User Does not have access to this route! Admin Only!",
				})
			}

			if config.RequiredAccess == AccessLevelUser && User.IsDeleted {
				return echo.NewHTTPError(http.StatusForbidden, map[string]any{
					"success": false,
					"message": "User does not exist! Deleted!",
				})
			}

			c.Set("UserSessionDetails", User)
			slog.Debug("User Session Context Applied")

			return next(c)
		}
	}
}
