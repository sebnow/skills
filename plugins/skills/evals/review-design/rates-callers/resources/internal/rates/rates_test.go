package rates

import (
	"errors"
	"strings"
	"testing"
)

func TestLoadKeepsTheLatestQuotePerPair(t *testing.T) {
	b, err := Load(strings.NewReader(`as_of,from,to,rate
2026-03-02T16:00:00Z,EUR,PLN,4.3
2026-03-03T16:00:00Z,EUR,PLN,4.312345
2026-03-01T16:00:00Z,EUR,PLN,4.25
`))
	if err != nil {
		t.Fatal(err)
	}
	q, err := b.Quote("EUR", "PLN")
	if err != nil || q.RateMicro != 4_312_345 || q.AsOf.Day() != 3 {
		t.Fatalf("Quote = %+v, %v", q, err)
	}
	if _, err := b.Quote("PLN", "EUR"); !errors.Is(err, ErrNoRate) {
		t.Fatalf("reverse pair: err = %v, want ErrNoRate", err)
	}
}

func TestLoadRefusesBadRows(t *testing.T) {
	for name, row := range map[string]string{
		"too many decimals":     "2026-03-02T16:00:00Z,EUR,PLN,4.1234567",
		"zero rate":             "2026-03-02T16:00:00Z,EUR,PLN,0",
		"negative rate":         "2026-03-02T16:00:00Z,EUR,PLN,-4.3",
		"trailing dot":          "2026-03-02T16:00:00Z,EUR,PLN,4.",
		"unsupported currency":  "2026-03-02T16:00:00Z,EUR,JPY,161.2",
		"same currency":         "2026-03-02T16:00:00Z,EUR,EUR,1",
		"time without timezone": "2026-03-02T16:00:00,EUR,PLN,4.3",
	} {
		if _, err := Load(strings.NewReader("as_of,from,to,rate\n" + row + "\n")); err == nil {
			t.Errorf("%s: Load accepted %q", name, row)
		}
	}
}
