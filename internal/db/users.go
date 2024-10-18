package db

import (
	"context"

	"github.com/HarshPatel5940/CodeFlick/internal/models"
	"github.com/jmoiron/sqlx"
)

const (
	getUserByID    = "SELECT * FROM users WHERE id = $1;"
	getUserByEmail = "SELECT * FROM users WHERE email = $1;"
	insertUser     = `INSERT INTO users (id, name, email, auth_provider, is_admin, is_premium, is_deleted)
		VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (email) DO NOTHING;`
)

type UserDB struct {
	db *sqlx.DB
}

func NewUserDB(db *sqlx.DB) *UserDB {
	return &UserDB{db: db}
}

func (u *UserDB) GetUserByID(ctx context.Context, userID string) (models.User, error) {
	var user models.User
	err := u.db.GetContext(ctx, &user, getUserByID, userID)
	return user, err
}

func (u *UserDB) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	err := u.db.GetContext(ctx, &user, getUserByEmail, email)
	return user, err
}

func (u *UserDB) InsertUser(ctx context.Context, user models.User) error {
	_, err := u.db.ExecContext(ctx, insertUser, user.ID, user.Name, user.Email, user.AuthProvider, user.IsAdmin, user.IsPremium, user.IsDeleted)
	return err
}
