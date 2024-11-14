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
	getAllUsers = `SELECT * FROM users LIMIT 500;`
)

type UserDB struct {
	db *sqlx.DB
	cm *RetryManager
}

func NewUserDB(db *sqlx.DB, cm *RetryManager) *UserDB {
	return &UserDB{db: db, cm: cm}
}

func (u *UserDB) GetUserByID(ctx context.Context, userID string) (models.User, error) {
	var user models.User
	err := u.cm.RetryWithSingleFlight(ctx, func() error {
		return u.db.GetContext(ctx, &user, getUserByID, userID)
	})
	return user, err
}

func (u *UserDB) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	err := u.cm.RetryWithSingleFlight(ctx, func() error {
		return u.db.GetContext(ctx, &user, getUserByEmail, email)
	})
	return user, err
}

func (u *UserDB) InsertUser(ctx context.Context, user models.User) error {
	return u.cm.RetryWithSingleFlight(ctx, func() error {
		_, err := u.db.ExecContext(ctx, insertUser, user.ID, user.Name, user.Email, user.AuthProvider, user.IsAdmin, user.IsPremium, user.IsDeleted)
		return err
	})
}

func (u *UserDB) GetAllUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User
	err := u.cm.RetryWithSingleFlight(ctx, func() error {
		return u.db.SelectContext(ctx, &users, getAllUsers)
	})
	return users, err
}
