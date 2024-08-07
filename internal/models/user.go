package models

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
