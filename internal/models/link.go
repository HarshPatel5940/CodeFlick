package models

import "time"

type File struct {
	UserID      string `json:"userId" db:"user_id"`
	ShortCode   string `json:"shortCode" db:"short_code"`
	OrgFileName string `json:"orgFileName" db:"org_file_name"`

	ViewCount int  `json:"viewCount" db:"view_count"`
	IsEnabled bool `json:"isEnabled" db:"is_enabled"`
	IsDeleted bool `json:"isDeleted" db:"is_deleted"`

	ExpiresAt time.Time `json:"expiresAt" db:"expires_at"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}
