package db

import (
	"context"
	"time"

	"github.com/HarshPatel5940/CodeFlick/internal/models"
	"github.com/jmoiron/sqlx"
)

const (
	getRepliesByGistID = `SELECT * FROM replies WHERE gist_id = $1;`
	getReplyByID       = `SELECT * FROM replies WHERE id = $1 AND gist_id = $2;`
	insertReply        = `INSERT INTO replies (id, user_id, gist_id, message) VALUES ($1, $2, $3, $4);`
	updateReply        = `UPDATE replies SET message = $4, updated_at = $5 WHERE id = $1 AND user_id = $2 AND gist_id = $3;`
	deleteReply        = `UPDATE replies SET is_deleted = true, updated_at = $4 WHERE id = $1 AND user_id = $2 AND gist_id = $3;`
)

type ReplyDB struct {
	db *sqlx.DB
}

func NewReplyDB(db *sqlx.DB) *ReplyDB {
	return &ReplyDB{db: db}
}

func (r *ReplyDB) GetRepliesByGistID(ctx context.Context, gistID string) ([]models.Reply, error) {
	var replies []models.Reply
	err := r.db.SelectContext(ctx, &replies, getRepliesByGistID, gistID)
	return replies, err
}

func (r *ReplyDB) GetReplyByID(ctx context.Context, replyID, gistID string) (models.Reply, error) {
	var reply models.Reply
	err := r.db.GetContext(ctx, &reply, getReplyByID, replyID, gistID)
	return reply, err
}

func (r *ReplyDB) InsertReply(ctx context.Context, reply models.Reply) error {
	_, err := r.db.ExecContext(ctx, insertReply, reply.ID, reply.UserID, reply.GistID, reply.Message)
	return err
}

func (r *ReplyDB) UpdateReply(ctx context.Context, reply models.Reply) error {
	_, err := r.db.ExecContext(ctx, updateReply, reply.ID, reply.UserID, reply.GistID, reply.Message, time.Now())
	return err
}

func (r *ReplyDB) DeleteReply(ctx context.Context, replyID, userID, gistID string) error {
	_, err := r.db.ExecContext(ctx, deleteReply, replyID, userID, gistID, time.Now())
	return err
}
