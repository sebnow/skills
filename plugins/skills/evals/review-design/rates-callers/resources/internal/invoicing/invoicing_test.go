package invoicing

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/brightloop/finops/internal/rates"
)

var now = time.Date(2026, 3, 3, 9, 0, 0, 0, time.UTC)

func book(t *testing.T) *rates.Book {
	t.Helper()
	b, err := rates.Load(strings.NewReader(`as_of,from,to,rate
2026-03-02T16:00:00Z,EUR,PLN,4.25
2026-03-02T16:00:00Z,USD,PLN,3.900005
2026-03-01T16:00:00Z,GBP,PLN,5.1
`))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestTotalConvertsAndRoundsHalfToEven(t *testing.T) {
	inv := Invoice{Customer: "acme", Currency: "PLN", Lines: []Line{
		{"setup", 10000, "PLN"},
		{"licence", 2, "EUR"},  // 8.5 grosz rounds to 8
		{"support", 50, "USD"}, // 195.00025 grosz rounds to 195
		{"credit", -2, "EUR"},  // -8.5 grosz rounds to -8
	}}
	got, err := Total(book(t), inv, now)
	if err != nil || got != 10000+8+195-8 {
		t.Fatalf("Total = %d, %v", got, err)
	}
}

func TestTotalRefusesStaleAndMissingRates(t *testing.T) {
	if _, err := Total(book(t), Invoice{Customer: "acme", Currency: "PLN", Lines: []Line{{"x", 100, "GBP"}}}, now); err == nil || !strings.Contains(err.Error(), "is stale") {
		t.Fatalf("stale rate: err = %v", err)
	}
	if _, err := Total(book(t), Invoice{Customer: "acme", Currency: "PLN", Lines: []Line{{"x", 100, "CHF"}}}, now); !errors.Is(err, rates.ErrNoRate) {
		t.Fatalf("missing rate: err = %v", err)
	}
}
