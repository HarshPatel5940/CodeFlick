package handlers

import (
	"context"
	"database/sql"
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
	sess, err := session.Get("session", c)

	if err != nil || sess.Values["email"] == nil {
		ctx := context.WithValue(c.Request().Context(), "provider", "google")
		gothic.BeginAuthHandler(c.Response(), c.Request().WithContext(ctx))
	}

	return c.JSON(200, map[string]any{
		"success": true,
		"message": "Session details fetched successfully!",
		"data": map[string]any{
			"name":      sess.Values["name"],
			"email":     sess.Values["email"],
			"isAdmin":   sess.Values["isAdmin"],
			"isPremium": sess.Values["isPremium"],
			"isCreated": sess.Values["created_at"],
		}})
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

	sess.Values["user_id"] = user.UserID
	sess.Values["name"] = user.Name
	sess.Values["email"] = user.Email
	sess.Values["auth_provider"] = user.Provider
	sess.Values["isAdmin"] = false
	sess.Values["isPremium"] = false
	sess.Values["isDeleted"] = false

	if err := ah.UpsertUserDetails(c, sess); err != nil {
		return err
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
			"name":      sess.Values["name"],
			"email":     sess.Values["email"],
			"isAdmin":   sess.Values["is_admin"],
			"isPremium": sess.Values["is_premium"],
			"isCreated": sess.Values["created_at"],
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

	query := `
		INSERT INTO users (id, name, email, auth_provider, is_admin, is_premium, is_deleted)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO NOTHING
		RETURNING id, created_at, updated_at;
		`

	row := Tx.QueryRowContext(
		context.Background(),
		query,
		sess.Values["user_id"],
		sess.Values["name"],
		sess.Values["email"],
		sess.Values["auth_provider"],
		sess.Values["isAdmin"],
		sess.Values["isPremium"],
		sess.Values["isDeleted"],
	)
	err = row.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			if err := Tx.QueryRowContext(
				context.Background(),
				"SELECT is_admin, is_premium, is_deleted, created_at FROM users WHERE id = $1",
				sess.Values["user_id"],
			).Scan(&user.IsAdmin, &user.IsPremium, &user.IsDeleted, &user.CreatedAt); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
					"success": false,
					"message": "Failed to fetch user details from PostgreSQL!",
					"details": err.Error(),
				})
			}
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
				"success": false,
				"message": "Failed to insert user details into PostgreSQL!",
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

	sess.Values["isAdmin"] = user.IsAdmin
	sess.Values["isPremium"] = user.IsPremium
	sess.Values["isDeleted"] = user.IsDeleted
	sess.Values["created_at"] = user.CreatedAt

	return nil
}
