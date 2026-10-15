package server

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
)

//go:embed schema.sql
var schema string

// migrations run in order after schema.sql; migration i takes the
// database from version i+1 to i+2.
var migrations = []string{
	`ALTER TABLE tasks ADD COLUMN retry_of TEXT REFERENCES tasks (id);`,
	`CREATE TABLE comments (
		id INTEGER PRIMARY KEY,
		task_id TEXT NOT NULL REFERENCES tasks (id),
		author TEXT NOT NULL,
		body TEXT NOT NULL,
		created TEXT NOT NULL
	);`,
	`CREATE TABLE artifacts (
		task_id TEXT NOT NULL REFERENCES tasks (id),
		name TEXT NOT NULL,
		size INTEGER NOT NULL,
		blob BLOB NOT NULL,
		PRIMARY KEY (task_id, name)
	);`,
	`CREATE TABLE requests (
		id INTEGER PRIMARY KEY,
		task_id TEXT NOT NULL REFERENCES tasks (id),
		question TEXT NOT NULL,
		answer TEXT
	);`,
	`CREATE TABLE task_links (
		from_task TEXT NOT NULL REFERENCES tasks (id),
		to_task TEXT NOT NULL REFERENCES tasks (id),
		kind TEXT NOT NULL
	);`,
	`CREATE TABLE stop_notes (
		task_id TEXT PRIMARY KEY REFERENCES tasks (id),
		note TEXT NOT NULL
	);`,
	`CREATE TABLE placements (
		task_id TEXT NOT NULL REFERENCES tasks (id),
		runner TEXT NOT NULL REFERENCES runners (id),
		placed TEXT NOT NULL
	);`,
	`CREATE TABLE notifications (
		id INTEGER PRIMARY KEY,
		task_id TEXT REFERENCES tasks (id),
		kind TEXT NOT NULL,
		read INTEGER NOT NULL DEFAULT 0
	);`,
}

// Store is the server's state in SQLite.
type Store struct {
	db *sql.DB
}

// migrate brings the database to the latest version.
func (s *Store) migrate(ctx context.Context) error {
	var version int
	if err := s.db.QueryRowContext(ctx, `PRAGMA user_version`).Scan(&version); err != nil {
		return fmt.Errorf("read schema version: %w", err)
	}
	if version == 0 {
		if _, err := s.db.ExecContext(ctx, schema); err != nil {
			return fmt.Errorf("create schema: %w", err)
		}
		version = 1
	}
	for i := version - 1; i < len(migrations); i++ {
		if _, err := s.db.ExecContext(ctx, migrations[i]); err != nil {
			return fmt.Errorf("migrate to version %d: %w", i+2, err)
		}
		if _, err := s.db.ExecContext(ctx, fmt.Sprintf(`PRAGMA user_version = %d`, i+2)); err != nil {
			return fmt.Errorf("set schema version %d: %w", i+2, err)
		}
	}
	return nil
}
