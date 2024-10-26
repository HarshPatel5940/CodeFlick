package models

import (
	"database/sql"
	"time"
)

type Gist struct {
	FileID string `json:"fileId" db:"file_id"`
	UserID string `json:"userId" db:"user_id"`

	GistTitle  string         `json:"gistTitle" db:"gist_title"`
	ForkedFrom sql.NullString `json:"forkedFrom" db:"forked_from"`
	ShortUrl   string         `json:"shortUrl" db:"short_url"`

	ViewCount int  `json:"viewCount" db:"view_count"`
	IsPublic  bool `json:"isPublic" db:"is_public"`
	IsDeleted bool `json:"isDeleted" db:"is_deleted"`

	AuditLog []string `json:"auditLog" db:"audit_logs"`

	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}
