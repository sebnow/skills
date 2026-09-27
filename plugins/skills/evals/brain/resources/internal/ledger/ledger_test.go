package ledger

import (
	"context"
	"path/filepath"
	"testing"
)

func TestOpenTwice(t *testing.T) {
	path := filepath.Join(t.TempDir(), "ledger.db")
	for i := 0; i < 2; i++ {
		lg, err := Open(context.Background(), path)
		if err != nil {
			t.Fatalf("open %d: %v", i, err)
		}
		lg.Close()
	}
}
