package store

import (
	"database/sql"
	"time"
)

// Session represents a chat/run session.
type Session struct {
	ID        int64
	CreatedAt time.Time
	Cwd       string
	RepoRoot  string
	Provider  string
	Model     string
	Title     string
}

// CreateSession inserts a new session and returns its ID.
func (s *Store) CreateSession(cwd, repoRoot, provider, model, title string) (int64, error) {
	res, err := s.DB.Exec(
		`INSERT INTO sessions (cwd, repo_root, provider, model, title) VALUES (?, ?, ?, ?, ?)`,
		cwd, repoRoot, provider, model, title,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// GetSession returns a session by ID.
func (s *Store) GetSession(id int64) (*Session, error) {
	var sess Session
	var createdAt string
	err := s.DB.QueryRow(
		`SELECT id, created_at, cwd, repo_root, provider, model, title FROM sessions WHERE id = ?`,
		id,
	).Scan(&sess.ID, &createdAt, &sess.Cwd, &sess.RepoRoot, &sess.Provider, &sess.Model, &sess.Title)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	sess.CreatedAt = parseTime(createdAt)
	return &sess, nil
}

// ListSessions returns recent sessions (e.g. last 50).
func (s *Store) ListSessions(limit int) ([]Session, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := s.DB.Query(
		`SELECT id, created_at, cwd, repo_root, provider, model, title FROM sessions ORDER BY id DESC LIMIT ?`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Session
	for rows.Next() {
		var sess Session
		var createdAt string
		if err := rows.Scan(&sess.ID, &createdAt, &sess.Cwd, &sess.RepoRoot, &sess.Provider, &sess.Model, &sess.Title); err != nil {
			return nil, err
		}
		sess.CreatedAt = parseTime(createdAt)
		out = append(out, sess)
	}
	return out, rows.Err()
}
