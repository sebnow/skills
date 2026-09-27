package ledger

import (
	"context"
	"time"
)

// Remove removes an import: its blob stays, a tombstone records the removal,
// and every election citing the blob as evidence goes with it (REQ-0104,
// REQ-0129).
func (l *Ledger) Remove(ctx context.Context, blob int64, reason string) error {
	tx, err := l.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx,
		`INSERT INTO tombstones (blob_id, removed_at, reason) VALUES (?, ?, ?)`,
		blob, l.now().UTC().Format(time.RFC3339Nano), reason); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `
		DELETE FROM election_prereqs WHERE election_id IN (SELECT election_id FROM election_evidence WHERE blob_id = ?)
		   OR requires IN (SELECT election_id FROM election_evidence WHERE blob_id = ?)`, blob, blob); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `
		DELETE FROM elections WHERE id IN (SELECT election_id FROM election_evidence WHERE blob_id = ?)`, blob); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM election_evidence WHERE blob_id = ?`, blob); err != nil {
		return err
	}
	return tx.Commit()
}
