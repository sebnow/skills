package export

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestBeancountBuy(t *testing.T) {
	var buf bytes.Buffer
	txs := []Transaction{{
		Date: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC), Custodian: "IBKR",
		Symbol: "VWRL", Quantity: 10, Price: 100, Currency: "EUR", Fee: 1.5,
	}}
	if err := Beancount(txs, &buf); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "Assets:IBKR:VWRL") {
		t.Fatalf("missing position posting:\n%s", buf.String())
	}
}
