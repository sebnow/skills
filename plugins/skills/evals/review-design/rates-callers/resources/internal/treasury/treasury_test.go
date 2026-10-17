package treasury

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

func TestPositionValuesBalancesAtTheLatestRates(t *testing.T) {
	got, err := Position(book(t), []Balance{{"main", 100000, "PLN"}, {"eu", 20000, "EUR"}}, "PLN", now)
	if err != nil || got != 100000+85000 {
		t.Fatalf("Position = %d, %v", got, err)
	}
}

func TestPositionAcceptsARateExactlyADayOld(t *testing.T) {
	if _, err := Position(book(t), []Balance{{"eu", 1, "EUR"}}, "PLN", time.Date(2026, 3, 3, 16, 0, 0, 0, time.UTC)); err != nil {
		t.Fatal(err)
	}
	_, err := Position(book(t), []Balance{{"eu", 1, "EUR"}}, "PLN", time.Date(2026, 3, 3, 16, 0, 1, 0, time.UTC))
	if err == nil || !strings.Contains(err.Error(), "is stale") || errors.Is(err, rates.ErrNoRate) {
		t.Fatalf("one second past a day: err = %v", err)
	}
}
