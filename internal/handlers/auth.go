package handlers

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/HarshPatel5940/CodeFlick/internal/utils"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

type AuthHandler struct {
	db         *sqlx.DB
	sessionAge int
}

func NewAuthHandler(db *sqlx.DB) *AuthHandler {
	InitialiseAuth()
	sessionAge, err := strconv.Atoi(utils.GetEnv("GORILLA_SESSIONS_MAXAGE", "604800"))

	if err != nil {
		slog.Error("GORILLA_SESSIONS_MAXAGE is not a valid integer! Taking '604800' (7 Days) as default value")
		slog.Error(err.Error())
		sessionAge = 604800
	}

	return &AuthHandler{db: db, sessionAge: sessionAge}
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

func (ah AuthHandler) GoogleOauthLogin(c echo.Context) error {

	// try to get the user without re-authenticating
	ctx := context.WithValue(c.Request().Context(), "provider", "google")
	gothUser, err := gothic.CompleteUserAuth(c.Response(), c.Request().WithContext(ctx))
	if err != nil {
		gothic.BeginAuthHandler(c.Response(), c.Request().WithContext(ctx))
	} else {
		log.Println(gothUser)
	}

	return nil
}

func (ah AuthHandler) GoogleOauthCallback(c echo.Context) error {
	user, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		return err
	}

	sess, err := session.Get("session", c)
	if err != nil {
		return err
	}

	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   ah.sessionAge,
		HttpOnly: true,
	}

	sess.Values["id"] = user.UserID
	sess.Values["name"] = user.Name
	sess.Values["email"] = user.Email
	sess.Values["verified_email"] = user.RawData["verified_email"]

	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed to save session!",
			"details": err.Error(),
		})

	}

	return c.JSON(200, map[string]any{
		"success": true,
		"message": "Successfully logged in with Google!",
	})
}

// Debug session details using this endpoint

func (ah AuthHandler) GetSessionDetails(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed to get session!",
			"details": err.Error(),
		})
	}

	return c.JSON(200, map[string]any{
		"success": true,
		"message": "Session details fetched successfully!",
		"data": map[string]any{
			"id":             sess.Values["id"],
			"name":           sess.Values["name"],
			"email":          sess.Values["email"],
			"verified_email": sess.Values["verified_email"]},
	})
}
