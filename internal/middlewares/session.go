package middlewares

import (
	"log/slog"
	"net/http"

	"github.com/HarshPatel5940/CodeFlick/internal/models"
	"github.com/HarshPatel5940/CodeFlick/internal/utils"
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

func SessionMiddleware(config Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			slog.Debug("Session Middleware Hit")
			sess, err := session.Get("session", c)

			if err != nil {
				if config.RequiredAccess != AccessLevelAll {
					return echo.NewHTTPError(http.StatusUnauthorized, map[string]any{
						"success": false,
						"message": "Failed to get session!",
					})
				}
			}

			userID := utils.GetSessionValue(sess, "user_id", "")
			if userID == "" && config.RequiredAccess != AccessLevelAll {
				return echo.NewHTTPError(http.StatusUnauthorized, map[string]any{
					"success": false,
					"message": "User not authenticated",
				})
			}

			User := models.User{
				ID:           userID,
				Name:         utils.GetSessionValue(sess, "name", ""),
				Email:        utils.GetSessionValue(sess, "email", ""),
				AuthProvider: utils.GetSessionValue(sess, "auth_provider", ""),
				IsAdmin:      utils.GetSessionValue(sess, "isAdmin", false),
				IsPremium:    utils.GetSessionValue(sess, "isPremium", false),
				IsDeleted:    utils.GetSessionValue(sess, "isDeleted", false),
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
			slog.Info("User Session Context Applited")

			return next(c)
		}
	}
}
