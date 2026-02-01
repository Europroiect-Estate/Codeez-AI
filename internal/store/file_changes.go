package store

import (
	"database/sql"
	"time"
)

// FileChange records a file change (diff) in a session.
type FileChange struct {
	ID        int64
	SessionID int64
	Path      string
	Diff      string
	Applied   bool
	CreatedAt time.Time
}

// RecordFileChange inserts a file change and returns its ID.
func (s *Store) RecordFileChange(sessionID int64, path, diff string) (int64, error) {
	res, err := s.DB.Exec(
		`INSERT INTO file_changes (session_id, path, diff) VALUES (?, ?, ?)`,
		sessionID, path, diff,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// MarkFileChangeApplied sets applied=1 for the given change ID.
func (s *Store) MarkFileChangeApplied(changeID int64) error {
	_, err := s.DB.Exec(`UPDATE file_changes SET applied = 1 WHERE id = ?`, changeID)
	return err
}

// GetFileChange returns a file change by ID.
func (s *Store) GetFileChange(id int64) (*FileChange, error) {
	var f FileChange
	var appliedInt int
	var createdAt string
	err := s.DB.QueryRow(
		`SELECT id, session_id, path, diff, applied, created_at FROM file_changes WHERE id = ?`,
		id,
	).Scan(&f.ID, &f.SessionID, &f.Path, &f.Diff, &appliedInt, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	f.Applied = appliedInt != 0
	f.CreatedAt = parseTime(createdAt)
	return &f, nil
}

// GetFileChangesForSession returns all file changes for a session.
func (s *Store) GetFileChangesForSession(sessionID int64) ([]FileChange, error) {
	rows, err := s.DB.Query(
		`SELECT id, session_id, path, diff, applied, created_at FROM file_changes WHERE session_id = ? ORDER BY id`,
		sessionID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []FileChange
	for rows.Next() {
		var f FileChange
		var appliedInt int
		var createdAt string
		if err := rows.Scan(&f.ID, &f.SessionID, &f.Path, &f.Diff, &appliedInt, &createdAt); err != nil {
			return nil, err
		}
		f.Applied = appliedInt != 0
		f.CreatedAt = parseTime(createdAt)
		out = append(out, f)
	}
	return out, rows.Err()
}
