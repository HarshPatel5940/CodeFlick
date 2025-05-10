package db

import (
	"context"
	"time"

	"github.com/HarshPatel5940/CodeFlick/internal/models"
	"github.com/jmoiron/sqlx"
)

/*
	type Reply struct {
		ID        string    `json:"id" db:"id"`
		UserID    string    `json:"user_id" db:"user_id"`
		GistID    string    `json:"gist_id" db:"gist_id"`
		Message   string    `json:"message" db:"message"`
		IsDeleted bool      `json:"is_deleted" db:"is_deleted"`
		CreatedAt time.Time `json:"created_at" db:"created_at"`
		UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	}

	type User struct {
		ID           string `json:"id" db:"id"`
		Name         string `json:"name" db:"name"`
		Email        string `json:"email" db:"email" `
		AuthProvider string `json:"authProvider" db:"auth_provider"`
		IsAdmin      bool   `json:"isAdmin" db:"is_admin"`
		IsPremium    bool   `json:"isPremium" db:"is_premium"`
		IsDeleted    bool   `json:"isDeleted" db:"is_deleted"`

		CreatedAt string `json:"createdAt" db:"created_at"`
		UpdatedAt string `json:"updatedAt" db:"updated_at"`
	}
*/
const (
	getRepliesByGistID = `
		SELECT r.*, u.name 
		FROM replies r 
		JOIN users u ON r.user_id = u.id 
		WHERE r.gist_id = $1;`
	getReplyByID = `
		SELECT r.*, u.name 
		FROM replies r 
		JOIN users u ON r.user_id = u.id 
		WHERE r.id = $1 AND r.gist_id = $2;`
	insertReply = `
		INSERT INTO replies (id, user_id, gist_id, message) 
		VALUES ($1, $2, $3, $4);`
	updateReply = `
		UPDATE replies 
		SET message = $1, updated_at = $2 
		WHERE id = $3 AND user_id = $4 AND gist_id = $5;`
	deleteReply = `
		UPDATE replies 
		SET is_deleted = true, updated_at = $1 
		WHERE id = $2 AND user_id = $3 AND gist_id = $4;`
)

type ReplyDB struct {
	db *sqlx.DB
	cm *RetryManager
}

func NewReplyDB(db *sqlx.DB, cm *RetryManager) *ReplyDB {
	return &ReplyDB{db: db, cm: cm}
}

func (r *ReplyDB) GetRepliesByGistID(ctx context.Context, gistID string) ([]models.Reply, error) {
	var replies []models.Reply
	err := r.cm.RetryWithSingleFlight(ctx, func() error {
		return r.db.SelectContext(ctx, &replies, getRepliesByGistID, gistID)
	})
	return replies, err
}

func (r *ReplyDB) GetReplyByID(ctx context.Context, replyID, gistID string) (models.Reply, error) {
	var reply models.Reply
	err := r.cm.RetryWithSingleFlight(ctx, func() error {
		return r.db.GetContext(ctx, &reply, getReplyByID, replyID, gistID)
	})
	return reply, err
}

func (r *ReplyDB) InsertReply(ctx context.Context, reply models.Reply) error {
	return r.cm.RetryWithSingleFlight(ctx, func() error {
		_, err := r.db.ExecContext(ctx, insertReply, reply.ID, reply.UserID, reply.GistID, reply.Message)
		return err
	})
}

func (r *ReplyDB) UpdateReply(ctx context.Context, reply models.Reply) error {
	return r.cm.RetryWithSingleFlight(ctx, func() error {
		_, err := r.db.ExecContext(ctx, updateReply, reply.ID, reply.UserID, reply.GistID, reply.Message, time.Now())
		return err
	})
}

func (r *ReplyDB) DeleteReply(ctx context.Context, replyID, userID, gistID string) error {
	return r.cm.RetryWithSingleFlight(ctx, func() error {
		_, err := r.db.ExecContext(ctx, deleteReply, replyID, userID, gistID, time.Now())
		return err
	})
}
