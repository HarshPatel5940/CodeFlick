package db

// KEEP queries constants here and then the functions that use them also here
const (
	GetUserByID    = "SELECT * FROM users WHERE id = $1;"
	GetUserByEmail = "SELECT * FROM users WHERE email = $1;"
	InsertUser     = `INSERT INTO users (id, name, email, auth_provider, is_admin, is_premium, is_deleted)
		VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (email) DO NOTHING;`

	GetGistByID       = `SELECT * FROM gists WHERE file_id = $1`
	GetGistByShortURL = `UPDATE gists SET view_count = view_count + 1 WHERE short_url = $1 RETURNING *;`
	GetGistsByUserID  = `SELECT * FROM gists WHERE user_id = $1;`
	InsertGist        = `INSERT INTO gists ( user_id, file_id, gist_title, short_url, is_public) VALUES ( $1, $2, $3, $4, $5 );`
	UpdateGist        = `UPDATE gists SET gist_title = $3, short_url = $4, is_public = $5, updated_at = $6 WHERE file_id = $1 and user_id=$2 RETURNING file_id;`
	DeleteGist        = `UPDATE gists SET is_deleted = true, updated_at = $3 WHERE file_id = $1 AND user_id=$2 RETURNING file_id;`
)
