package flow

import (
	"context"
	"path/filepath"
	"testing"

	"tidewright/internal/ledger"
)

func TestElectionsInForceSkipsSuperseded(t *testing.T) {
	ctx := context.Background()
	lg, err := ledger.Open(ctx, filepath.Join(t.TempDir(), "ledger.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer lg.Close()
	if _, err := lg.Elect(ctx, "levy-instalments", "single", nil); err != nil {
		t.Fatal(err)
	}
	if _, err := lg.Elect(ctx, "levy-instalments", "two", nil); err != nil {
		t.Fatal(err)
	}
	es, err := electionsInForce(ctx, lg)
	if err != nil {
		t.Fatal(err)
	}
	if len(es) != 1 || es[0].Option != "two" {
		t.Fatalf("elections in force = %+v, want only levy-instalments=two", es)
	}
}
