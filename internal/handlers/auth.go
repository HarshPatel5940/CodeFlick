package handlers

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/HarshPatel5940/CodeFlick/internal/models"
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
		Secure:   true,
		MaxAge:   ah.sessionAge,
		HttpOnly: true,
	}

	sess.Values["name"] = user.Name
	sess.Values["email"] = user.Email
	sess.Values["auth_provider"] = user.Provider
	sess.Values["is_admin"] = false
	sess.Values["is_premium"] = false
	sess.Values["is_deleted"] = false

	if err := ah.UpsertUserDetails(c, sess); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed to fetch user details!",
			"details": err.Error(),
		})
	}

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
			"name":          sess.Values["name"],
			"email":         sess.Values["email"],
			"auth_provider": sess.Values["auth_provider"],
			"is_admin":      sess.Values["is_admin"],
			"is_premium":    sess.Values["is_premium"],
		}})
}

func (ah AuthHandler) UpsertUserDetails(c echo.Context, sess *sessions.Session) error {
	Tx, err := ah.db.BeginTx(context.Background(), &sql.TxOptions{})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed to start a PostgreSQL transaction!",
			"details": err.Error(),
		})
	}

	defer Tx.Rollback()

	var user models.User

	if err := Tx.QueryRowContext(
		context.Background(),
		"SELECT * FROM users WHERE email = $1",
		sess.Values["email"],
	).Scan(&user.Name, &user.Email, &user.AuthProvider, &user.IsAdmin, &user.IsPremium, &user.IsDeleted); err != nil {
		if err == sql.ErrNoRows {
			res, err := Tx.Exec(
				"INSERT INTO users (name, email, auth_provider, is_admin, is_premium, is_deleted) VALUES ($1, $2, $3, $4, $5, $6)",
				sess.Values["name"],
				sess.Values["email"],
				sess.Values["auth_provider"],
				sess.Values["is_admin"],
				sess.Values["is_premium"],
				sess.Values["is_deleted"])

			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
					"success": false,
					"message": "Failed to insert user details into PostgreSQL!",
					"details": err.Error(),
				})

			}

			log.Println(res)
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
				"success": false,
				"message": "Failed to fetch user details from PostgreSQL!",
				"details": err.Error(),
			})
		}
	}

	if err := Tx.Commit(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed to commit the PostgreSQL transaction!",
			"details": err.Error(),
		})
	}

	sess.Values["is_admin"] = user.IsAdmin
	sess.Values["is_premium"] = user.IsPremium
	sess.Values["is_deleted"] = user.IsDeleted

	return nil

}
