package models

import "time"

type Reply struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	GistID    string    `json:"gist_id" db:"gist_id"`
	Message   string    `json:"message" db:"message"`
	IsDeleted bool      `json:"is_deleted" db:"is_deleted"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Name string `json:"name" db:"name"`
}
