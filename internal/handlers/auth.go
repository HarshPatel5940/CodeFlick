package handlers

import (
	"context"
	"log"

	"github.com/HarshPatel5940/CodeFlick/internal/utils"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

type AuthHandler struct {
	db *sqlx.DB
}

func NewAuthHandler(db *sqlx.DB) *AuthHandler {
	InitialiseAuth()
	return &AuthHandler{db: db}
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
			"email", "profile"),
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

/*
 // google user response
{
  "RawData": {
    "email": "REDACTED",
    "family_name": "Patel",
    "given_name": "Harsh",
    "id": "REDACTED",
    "name": "Harsh Patel",
    "picture": "REDACTED",
    "verified_email": true
  },
  "Provider": "google",
  "Email": "REDACTED",
  "Name": "Harsh Patel",
  "FirstName": "Harsh",
  "LastName": "Patel",
  "NickName": "Harsh Patel",
  "Description": "",
  "UserID": "REDACTED",
  "AvatarURL": "REDACTED"
  "Location": "",
  "AccessToken": "REDACTED",
  "AccessTokenSecret": "",
  "RefreshToken": "REDACTED",
  "ExpiresAt": "2024-07-02T11:29:03.613039+05:30",
  "IDToken": "REDACTED",
}
*/

func (ah AuthHandler) GoogleOauthCallback(c echo.Context) error {
	user, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		return err
	}
	// TODO: store the user in a session

	return c.JSON(200, user)
}
