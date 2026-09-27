package ledger

import (
	"context"
	"path/filepath"
	"testing"
)

func open(t *testing.T) *Ledger {
	t.Helper()
	lg, err := Open(context.Background(), filepath.Join(t.TempDir(), "ledger.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { lg.Close() })
	return lg
}

func TestElectSupersedes(t *testing.T) {
	ctx := context.Background()
	lg := open(t)
	first, err := lg.Elect(ctx, "survey-deferral", "none", nil)
	if err != nil {
		t.Fatal(err)
	}
	second, err := lg.Elect(ctx, "survey-deferral", "three-months", nil)
	if err != nil {
		t.Fatal(err)
	}
	if second.Supersedes != first.ID {
		t.Fatalf("second election supersedes %d, want %d", second.Supersedes, first.ID)
	}
	all, err := lg.Elections(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 2 {
		t.Fatalf("%d elections on record, want 2 (REQ-0140)", len(all))
	}
}
