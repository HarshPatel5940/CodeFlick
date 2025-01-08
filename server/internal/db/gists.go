package db

import (
	"context"
	"time"

	"github.com/HarshPatel5940/CodeFlick/internal/models"
	"github.com/jmoiron/sqlx"
)

const (
	getGistByID       = `SELECT * FROM gists WHERE file_id = $1`
	getGistByShortURL = `UPDATE gists SET view_count = view_count + 1 WHERE short_url = $1 RETURNING *;`
	getGistsByUserID  = `SELECT * FROM gists WHERE user_id = $1;`
	getGistsByUserIDOrPublic = `SELECT * FROM gists WHERE user_id = $1 OR is_public = true;`
	getPublicGists   = `SELECT * FROM gists WHERE is_public = true;`
	insertGist        = `INSERT INTO gists (user_id, file_id, gist_title, short_url, is_public) VALUES ($1, $2, $3, $4, $5);`
	updateGist        = `UPDATE gists SET gist_title = $3, short_url = $4, is_public = $5, updated_at = $6 WHERE file_id = $1 and user_id=$2 RETURNING file_id;`
	deleteGist        = `UPDATE gists SET is_deleted = true, updated_at = $3 WHERE file_id = $1 AND user_id=$2 RETURNING file_id;`
)

type GistDB struct {
	db *sqlx.DB
	cm *RetryManager
}

func NewGistDB(db *sqlx.DB, cm *RetryManager) *GistDB {
	return &GistDB{db: db, cm: cm}
}

func (g *GistDB) GetGistByID(ctx context.Context, gistID string) (models.Gist, error) {
	var gist models.Gist
	err := g.cm.RetryWithSingleFlight(ctx, func() error {
		return g.db.GetContext(ctx, &gist, getGistByID, gistID)
	})
	return gist, err
}

func (g *GistDB) GetGistByShortURL(ctx context.Context, shortURL string) (models.Gist, error) {
	var gist models.Gist
	err := g.cm.RetryWithSingleFlight(ctx, func() error {
		return g.db.GetContext(ctx, &gist, getGistByShortURL, shortURL)
	})
	return gist, err
}

func (g *GistDB) GetGistsByUserID(ctx context.Context, userID string, fetchPublic string) ([]models.Gist, error) {
	var gists []models.Gist = []models.Gist{}

	switch fetchPublic {
		case "yes":
			err := g.cm.RetryWithSingleFlight(ctx, func() error {
				return g.db.SelectContext(ctx, &gists, getGistsByUserIDOrPublic, userID)
			})
			return gists, err

		case "only":
			err := g.cm.RetryWithSingleFlight(ctx, func() error {
				return g.db.SelectContext(ctx, &gists, getPublicGists)
			})
			return gists, err

		default:
		err := g.cm.RetryWithSingleFlight(ctx, func() error {
				return g.db.SelectContext(ctx, &gists, getGistsByUserID, userID)
			})
			return gists, err
	}

}

func (g *GistDB) InsertGist(ctx context.Context, gist models.Gist) error {
	return g.cm.RetryWithSingleFlight(ctx, func() error {
		_, err := g.db.ExecContext(ctx, insertGist, gist.UserID, gist.FileID, gist.GistTitle, gist.ShortUrl, gist.IsPublic)
		return err
	})
}

func (g *GistDB) UpdateGist(ctx context.Context, gist models.Gist) (string, error) {
	var fileID string
	err := g.cm.RetryWithSingleFlight(ctx, func() error {
		return g.db.GetContext(ctx, &fileID, updateGist, gist.FileID, gist.UserID, gist.GistTitle, gist.ShortUrl, gist.IsPublic, time.Now())
	})
	return fileID, err
}

func (g *GistDB) DeleteGist(ctx context.Context, gistID, userID string) (string, error) {
	var fileID string
	err := g.cm.RetryWithSingleFlight(ctx, func() error {
		return g.db.GetContext(ctx, &fileID, deleteGist, gistID, userID, time.Now())
	})
	return fileID, err
}
