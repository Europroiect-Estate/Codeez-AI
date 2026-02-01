package store

import (
	"database/sql"
	"fmt"
)

const schemaVersion = 1

// Migrate runs idempotent migrations.
func Migrate(s *Store) error {
	if err := ensureVersionTable(s.DB); err != nil {
		return err
	}
	v, err := getVersion(s.DB)
	if err != nil {
		return err
	}
	if v >= schemaVersion {
		return nil
	}
	for _, stmt := range migrations {
		if _, err := s.DB.Exec(stmt); err != nil {
			return fmt.Errorf("migration: %w", err)
		}
	}
	return setVersion(s.DB, schemaVersion)
}

func ensureVersionTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_version (version INTEGER NOT NULL);
		INSERT OR IGNORE INTO schema_version (version) VALUES (0);
	`)
	return err
}

func getVersion(db *sql.DB) (int, error) {
	var v int
	err := db.QueryRow("SELECT version FROM schema_version LIMIT 1").Scan(&v)
	return v, err
}

func setVersion(db *sql.DB, version int) error {
	_, err := db.Exec("UPDATE schema_version SET version = ?", version)
	return err
}

var migrations = []string{
	`CREATE TABLE IF NOT EXISTS sessions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at TEXT NOT NULL DEFAULT (datetime('now')),
		cwd TEXT NOT NULL,
		repo_root TEXT,
		provider TEXT,
		model TEXT,
		title TEXT
	);`,
	`CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		session_id INTEGER NOT NULL REFERENCES sessions(id),
		role TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at TEXT NOT NULL DEFAULT (datetime('now'))
	);`,
	`CREATE TABLE IF NOT EXISTS tool_audit (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		session_id INTEGER REFERENCES sessions(id),
		agent TEXT NOT NULL,
		tool_name TEXT NOT NULL,
		args_json TEXT,
		approved INTEGER NOT NULL,
		approval_scope TEXT,
		created_at TEXT NOT NULL DEFAULT (datetime('now'))
	);`,
	`CREATE TABLE IF NOT EXISTS file_changes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		session_id INTEGER REFERENCES sessions(id),
		path TEXT NOT NULL,
		diff TEXT,
		applied INTEGER NOT NULL DEFAULT 0,
		created_at TEXT NOT NULL DEFAULT (datetime('now'))
	);`,
	`CREATE TABLE IF NOT EXISTS permissions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		repo_root TEXT NOT NULL,
		tool_name TEXT NOT NULL,
		scope TEXT NOT NULL,
		rule_json TEXT,
		created_at TEXT NOT NULL DEFAULT (datetime('now'))
	);`,
}
