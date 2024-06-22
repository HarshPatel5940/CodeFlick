package models

type User struct {
	UserID       string `json:"userId" db:"user_id"`
	Name         string `json:"name" db:"name"`
	Email        string `json:"email" db:"email"`
	AuthProvider string `json:"authProvider" db:"auth_provider"`
	IsAdmin      bool   `json:"isAdmin" db:"is_admin"`
	IsDeleted    bool   `json:"isDeleted" db:"is_deleted"`

	CreatedAt string `json:"createdAt" db:"created_at"`
	UpdatedAt string `json:"updatedAt" db:"updated_at"`
}
