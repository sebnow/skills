package ledger

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// ElectionKind names what an election records. The kind decides how Body
// is decoded and which rule reads it.
type ElectionKind string

const (
	// KindOption records which option of a clause the operator
	// picked (REQ-0139).
	KindOption ElectionKind = "option"
)

func (k ElectionKind) valid() bool {
	switch k {
	case KindOption:
		return true
	default:
		return false
	}
}

// Election is one recorded election. About identifies what it is about;
// a later election with the same kind and about supersedes it (REQ-0140).
type Election struct {
	ID           int64
	Kind         ElectionKind
	About      string
	MadeAt   time.Time
	Body         json.RawMessage
	Sources      []Source
	Requires     []int64
	Supersedes   int64
	SupersededBy int64
}

// Source is a stored blob that supports an election, optionally narrowed to
// one record inside it.
type Source struct {
	Blob    int64
	Locator string
}

const electionSchema = `
CREATE TABLE IF NOT EXISTS elections (
	id          INTEGER PRIMARY KEY,
	kind        TEXT NOT NULL,
	about       TEXT NOT NULL,
	made_at     TEXT NOT NULL,
	body        TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS election_evidence (
	election_id INTEGER NOT NULL REFERENCES elections(id),
	blob_id     INTEGER NOT NULL REFERENCES blobs(id),
	locator     TEXT NOT NULL DEFAULT ''
);
CREATE TABLE IF NOT EXISTS election_prereqs (
	election_id INTEGER NOT NULL REFERENCES elections(id),
	requires    INTEGER NOT NULL REFERENCES elections(id)
);
CREATE INDEX IF NOT EXISTS elections_about ON elections(kind, about);
`

// OptionBody is the Body of a KindOption election.
type OptionBody struct {
	Clause string `json:"clause"`
	Option string `json:"option"`
}

// Elect stores the operator's option for clause, superseding any earlier
// election for the same clause.
func (l *Ledger) Elect(ctx context.Context, clause, option string, sources []Source) (Election, error) {
	body, err := json.Marshal(OptionBody{Clause: clause, Option: option})
	if err != nil {
		return Election{}, err
	}
	return l.put(ctx, KindOption, clause, body, sources, nil)
}

// put inserts an election of kind about its subject inside one transaction and
// links it to its evidence and prerequisites. Removal of a cited import
// cascades through election_evidence (REQ-0129).
func (l *Ledger) put(ctx context.Context, kind ElectionKind, about string, body []byte, sources []Source, requires []int64) (Election, error) {
	if !kind.valid() {
		return Election{}, fmt.Errorf("ledger: unknown election kind %q", kind)
	}
	tx, err := l.db.BeginTx(ctx, nil)
	if err != nil {
		return Election{}, err
	}
	defer tx.Rollback()
	now := l.now().UTC()
	res, err := tx.ExecContext(ctx,
		`INSERT INTO elections (kind, about, made_at, body) VALUES (?, ?, ?, ?)`,
		kind, about, now.Format(time.RFC3339Nano), string(body))
	if err != nil {
		return Election{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return Election{}, err
	}
	for _, s := range sources {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO election_evidence (election_id, blob_id, locator) VALUES (?, ?, ?)`,
			id, s.Blob, s.Locator); err != nil {
			return Election{}, err
		}
	}
	for _, r := range requires {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO election_prereqs (election_id, requires) VALUES (?, ?)`, id, r); err != nil {
			return Election{}, err
		}
	}
	var prev int64
	err = tx.QueryRowContext(ctx,
		`SELECT id FROM elections WHERE kind = ? AND about = ? AND id < ? ORDER BY id DESC LIMIT 1`,
		kind, about, id).Scan(&prev)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return Election{}, err
	}
	if err := tx.Commit(); err != nil {
		return Election{}, err
	}
	return Election{ID: id, Kind: kind, About: about, MadeAt: now, Body: body, Sources: sources, Requires: requires, Supersedes: prev}, nil
}

// Elections returns every election with its supersession resolved, newest
// first within each about.
func (l *Ledger) Elections(ctx context.Context) ([]Election, error) {
	rows, err := l.db.QueryContext(ctx, `
		SELECT e.id, e.kind, e.about, e.made_at, e.body,
		       COALESCE((SELECT MIN(n.id) FROM elections n
		                 WHERE n.kind = e.kind AND n.about = e.about AND n.id > e.id), 0)
		FROM elections e ORDER BY e.kind, e.about, e.id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Election
	for rows.Next() {
		var e Election
		var at, body string
		if err := rows.Scan(&e.ID, &e.Kind, &e.About, &at, &body, &e.SupersededBy); err != nil {
			return nil, err
		}
		e.MadeAt, _ = time.Parse(time.RFC3339Nano, at)
		e.Body = json.RawMessage(body)
		out = append(out, e)
	}
	return out, rows.Err()
}

// Summary renders the election's value for listings.
func (e Election) Summary() string {
	switch e.Kind {
	case KindOption:
		var b OptionBody
		if json.Unmarshal(e.Body, &b) == nil {
			return b.Option
		}
	}
	return string(e.Body)
}
