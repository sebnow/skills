package payouts

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

func TestAmountPaysInThePartnerCurrency(t *testing.T) {
	p := Payout{Partner: "kiosk-7", Currency: "PLN", Earnings: []Earning{{"o1", 1000, "EUR"}, {"o2", 300, "PLN"}}}
	got, err := Amount(book(t), p, now)
	if err != nil || got != 4250+300 {
		t.Fatalf("Amount = %d, %v", got, err)
	}
}

func TestAmountRefusesAStaleRate(t *testing.T) {
	p := Payout{Partner: "kiosk-7", Currency: "PLN", Earnings: []Earning{{"o1", 1000, "GBP"}}}
	if _, err := Amount(book(t), p, now); err == nil || !strings.Contains(err.Error(), "is stale") {
		t.Fatalf("err = %v", err)
	}
	if _, err := Amount(book(t), Payout{Partner: "kiosk-7", Currency: "PLN", Earnings: []Earning{{"o1", 1, "SEK"}}}, now); !errors.Is(err, rates.ErrNoRate) {
		t.Fatalf("missing rate: err = %v", err)
	}
}
