package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/HarshPatel5940/CodeFlick/internal/db"
	"github.com/HarshPatel5940/CodeFlick/internal/models"
	"github.com/HarshPatel5940/CodeFlick/internal/utils"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/oklog/ulid/v2"
)

type AuthHandler struct {
	userDB     *db.UserDB
	sessionAge int
	clientURL  string
}

func NewAuthHandler(userDB *db.UserDB) *AuthHandler {
	InitialiseAuth()
	sessionAge, err := strconv.Atoi(utils.GetEnv("GORILLA_SESSIONS_MAXAGE", "604800"))
	if err != nil {
		slog.Error("GORILLA_SESSIONS_MAXAGE is not a valid integer! Taking '604800' (7 Days) as default value")
		slog.Error(err.Error())
		sessionAge = 604800
	}

	clientURL := utils.GetEnv("CLIENT_URL")

	return &AuthHandler{userDB: userDB, sessionAge: sessionAge, clientURL: clientURL}
}

func InitialiseAuth() {
	store := sessions.NewCookieStore([]byte(
		utils.GetEnv("GORILLA_SESSIONS_KEY"),
	))

	gothic.Store = store

	goth.UseProviders(
		google.New(
			utils.GetEnv("GOOGLE_CLIENT_ID"),
			utils.GetEnv("GOOGLE_CLIENT_SECRET"),
			utils.GetEnv("GOOGLE_CALLBACK_URL"),
		),
	)
}

func (ah *AuthHandler) GoogleOauthLogin(c echo.Context) error {
	sess, err := session.Get("_gothic_temp_session", c)
	sess.Options = &sessions.Options{
		Path:     "/",
		Secure:   true,
		MaxAge:   ah.sessionAge,
		HttpOnly: true,
	}

	redirectParam := c.QueryParam("r")

	if err != nil || sess.Values["user_id"] == nil || sess.Values["user_id"] == "" {
		if redirectParam != "" {
			sess.Values["oauth_redirect"] = redirectParam
			if err := sess.Save(c.Request(), c.Response()); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save session")
			}
		}

		// ignore the below linting warning... i am having skill issue while fixing that :(
		ctx := context.WithValue(c.Request().Context(), "provider", "google")
		c.SetRequest(c.Request().WithContext(ctx))

		gothic.BeginAuthHandler(c.Response(), c.Request())
		return nil
	}

	if redirectParam == "self" {
		return c.Redirect(http.StatusFound, "/")
	}

	if redirectParam == "client" {
		return c.Redirect(http.StatusFound, ah.clientURL)
	}

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": "Session details fetched successfully! without callback",
		"data": map[string]any{
			"user_id":   sess.Values["user_id"],
			"name":      sess.Values["name"],
			"email":     sess.Values["email"],
			"isAdmin":   sess.Values["isAdmin"],
			"isPremium": sess.Values["isPremium"],
			"isDeleted": sess.Values["isDeleted"],
		},
	})
}

func (ah *AuthHandler) GoogleOauthCallback(c echo.Context) error {
	GothUser, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to complete user authentication")
	}

	sess, err := session.Get("_gothic_session", c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, map[string]any{
			"success": false,
			"message": "Failed to get session!",
		})
	}

	sess.Options = &sessions.Options{
		Path:     "/",
		Secure:   true,
		MaxAge:   ah.sessionAge,
		HttpOnly: true,
	}

	User, err := ah.userDB.GetUserByEmail(context.Background(), GothUser.Email)
	if err != nil {
		User = models.User{
			ID:           ulid.Make().String(),
			Name:         GothUser.Name,
			Email:        GothUser.Email,
			AuthProvider: GothUser.Provider,
			IsAdmin:      false,
			IsPremium:    false,
			IsDeleted:    false,
		}
		err = ah.userDB.InsertUser(context.Background(), User)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
				"success": false,
				"message": "Failed to insert user details!",
				"details": err.Error(),
			})
		}
	}

	sess.Values["user_id"] = User.ID
	sess.Values["name"] = User.Name
	sess.Values["email"] = User.Email
	sess.Values["auth_provider"] = User.AuthProvider
	sess.Values["isAdmin"] = User.IsAdmin
	sess.Values["isPremium"] = User.IsPremium
	sess.Values["isDeleted"] = User.IsDeleted

	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed to save session!",
		})
	}
	Tsess, err := session.Get("_gothic_temp_session", c)
	if err != nil {
		slog.Error(err.Error())
		return c.JSON(http.StatusOK, map[string]any{
			"success": true,
			"message": "Successfully logged in with Google!  Redirect Checking Failed...",
		})
	}

	redirectParam, ok := Tsess.Values["oauth_redirect"].(string)
	if !ok || redirectParam == "" {
		redirectParam = c.QueryParam("r")
	}

	Tsess.Values = make(map[any]any)
	delete(Tsess.Values, "oauth_redirect")
	Tsess.Options = &sessions.Options{
		MaxAge: -1,
	}
	if err := Tsess.Save(c.Request(), c.Response()); err != nil {
		slog.Error("Failed to clear oauth_redirect from session", "error", err)
	}

	if redirectParam == "self" {
		return c.Redirect(http.StatusFound, "/")
	}

	if redirectParam == "client" {
		return c.Redirect(http.StatusFound, ah.clientURL)
	}

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": "Successfully logged in with Google!",
	})
}

func (ah *AuthHandler) GetSessionDetails(c echo.Context) error {
	var sess models.User = c.Get("UserSessionDetails").(models.User)

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": "Session details fetched successfully!",
		"data": map[string]any{
			"user_id":   sess.ID,
			"name":      sess.Name,
			"email":     sess.Email,
			"isAdmin":   sess.IsAdmin,
			"isPremium": sess.IsPremium,
			"isDeleted": sess.IsDeleted,
		},
	})
}

func (ah *AuthHandler) Logout(c echo.Context) error {
	if err := gothic.Logout(c.Response(), c.Request()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Error during logout!",
			"details": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": "Successfully logged out!",
	})
}

func (ah *AuthHandler) GetAllUsers(c echo.Context) error {
	users, err := ah.userDB.GetAllUsers(context.Background())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed while fetching users from databases!",
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": "All Users details fetched successfully!",
		"data":    users,
	})
}
