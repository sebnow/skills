package ledger

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

// Ledger is the SQLite file that holds every input (ADR-0004). Imports go in
// as blobs and are never modified (ADR-0003); nothing derived is stored.
type Ledger struct {
	db  *sql.DB
	now func() time.Time
}

// Open opens or creates the ledger at path and applies the schema.
func Open(ctx context.Context, path string) (*Ledger, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	l := &Ledger{db: db, now: time.Now}
	for _, s := range []string{schema, electionSchema, recordSchema} {
		if _, err := db.ExecContext(ctx, s); err != nil {
			db.Close()
			return nil, fmt.Errorf("ledger: schema: %w", err)
		}
	}
	return l, nil
}

// Close closes the ledger file.
func (l *Ledger) Close() error { return l.db.Close() }

// schema holds the imported files. A removed import keeps its blob and gets a
// tombstone, so the removal can be explained (REQ-0104).
const schema = `
CREATE TABLE IF NOT EXISTS blobs (
	id      INTEGER PRIMARY KEY,
	name    TEXT NOT NULL,
	sha256  TEXT NOT NULL UNIQUE,
	size    INTEGER NOT NULL,
	stored_at TEXT NOT NULL,
	content BLOB NOT NULL
);
CREATE TABLE IF NOT EXISTS imports (
	blob_id     INTEGER NOT NULL REFERENCES blobs(id),
	name        TEXT NOT NULL,
	format      TEXT NOT NULL,
	records     INTEGER NOT NULL DEFAULT 0,
	imported_at TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS tombstones (
	blob_id    INTEGER NOT NULL REFERENCES blobs(id),
	removed_at TEXT NOT NULL,
	reason     TEXT NOT NULL DEFAULT ''
);
-- A blob is referenced by at most one import; names are unique so
-- `-evidence FILE` can find a stored document by name.
CREATE UNIQUE INDEX IF NOT EXISTS imports_blob ON imports(blob_id);
CREATE UNIQUE INDEX IF NOT EXISTS blobs_name ON blobs(name);
CREATE INDEX IF NOT EXISTS tombstones_blob ON tombstones(blob_id);
`

// recordSchema holds overhaul and survey records, which the operator enters
// with `tidewright overhaul` and `tidewright survey`.
const recordSchema = `
CREATE TABLE IF NOT EXISTS overhauls (engine TEXT NOT NULL, at TEXT NOT NULL);
CREATE TABLE IF NOT EXISTS surveys (at TEXT NOT NULL, certificate TEXT NOT NULL);
`

// Blob is a stored file.
type Blob struct {
	ID   int64
	Name string
}

// BlobByName returns the stored file with the given name.
func (l *Ledger) BlobByName(ctx context.Context, name string) (Blob, error) {
	var b Blob
	err := l.db.QueryRowContext(ctx, `SELECT id, name FROM blobs WHERE name = ?`, name).Scan(&b.ID, &b.Name)
	if err == sql.ErrNoRows {
		return Blob{}, fmt.Errorf("ledger: no stored file named %q", name)
	}
	return b, err
}

// Import is one imported file.
type Import struct {
	Blob       int64
	Name       string
	Format     string
	ImportedAt time.Time
	Removed    bool
}

// Imports lists every import, removed ones included, oldest first.
func (l *Ledger) Imports(ctx context.Context) ([]Import, error) {
	rows, err := l.db.QueryContext(ctx, `
		SELECT i.blob_id, i.name, i.format, i.imported_at,
		       EXISTS (SELECT 1 FROM tombstones t WHERE t.blob_id = i.blob_id)
		FROM imports i ORDER BY i.imported_at`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Import
	for rows.Next() {
		var im Import
		var at string
		if err := rows.Scan(&im.Blob, &im.Name, &im.Format, &at, &im.Removed); err != nil {
			return nil, err
		}
		im.ImportedAt, _ = time.Parse(time.RFC3339Nano, at)
		out = append(out, im)
	}
	return out, rows.Err()
}

// Content returns the bytes of a stored file.
func (l *Ledger) Content(ctx context.Context, blob int64) ([]byte, error) {
	var b []byte
	err := l.db.QueryRowContext(ctx, `SELECT content FROM blobs WHERE id = ?`, blob).Scan(&b)
	return b, err
}

// OverhaulRecord is an engine overhaul the operator recorded.
type OverhaulRecord struct {
	Engine string
	At     time.Time
}

// Overhauls lists recorded overhauls.
func (l *Ledger) Overhauls(ctx context.Context) ([]OverhaulRecord, error) {
	rows, err := l.db.QueryContext(ctx, `SELECT engine, at FROM overhauls ORDER BY at`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []OverhaulRecord
	for rows.Next() {
		var o OverhaulRecord
		var at string
		if err := rows.Scan(&o.Engine, &at); err != nil {
			return nil, err
		}
		o.At, _ = time.Parse("2006-01-02", at)
		out = append(out, o)
	}
	return out, rows.Err()
}

// LastSurvey returns the date on the latest survey certificate.
func (l *Ledger) LastSurvey(ctx context.Context) (time.Time, error) {
	var at sql.NullString
	if err := l.db.QueryRowContext(ctx, `SELECT MAX(at) FROM surveys`).Scan(&at); err != nil {
		return time.Time{}, err
	}
	if !at.Valid {
		return time.Time{}, nil
	}
	return time.Parse("2006-01-02", at.String)
}
