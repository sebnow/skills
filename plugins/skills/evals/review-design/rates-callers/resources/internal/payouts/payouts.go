// Package payouts works out what partners are paid for their orders.
package payouts

import (
	"fmt"
	"math/big"
	"time"

	"github.com/brightloop/finops/internal/rates"
)

type Earning struct {
	Order       string
	AmountCents int64
	Currency    string
}

type Payout struct {
	Partner  string
	Currency string
	Earnings []Earning
}

// Amount returns what the partner is paid, in the partner's currency.
func Amount(book *rates.Book, p Payout, now time.Time) (int64, error) {
	var total int64
	for _, e := range p.Earnings {
		amount, err := toPartnerCurrency(book, e.AmountCents, e.Currency, p.Currency, now)
		if err != nil {
			return 0, fmt.Errorf("payout to %s: %w", p.Partner, err)
		}
		total += amount
	}
	return total, nil
}

// toPartnerCurrency turns cents in one currency into cents in another at the latest rate,
// refusing a rate published more than a day before now.
func toPartnerCurrency(book *rates.Book, cents int64, from, to string, now time.Time) (int64, error) {
	if from == to {
		return cents, nil
	}
	q, err := book.Quote(from, to)
	if err != nil {
		return 0, err
	}
	if now.Sub(q.AsOf) > 24*time.Hour {
		return 0, fmt.Errorf("%s/%s rate from %s is stale", from, to, q.AsOf.Format(time.RFC3339))
	}
	// Round half to even, to the cent.
	p := new(big.Int).Mul(big.NewInt(cents), big.NewInt(q.RateMicro))
	quo, rem := new(big.Int).QuoRem(p, big.NewInt(1_000_000), new(big.Int))
	twice := new(big.Int).Lsh(new(big.Int).Abs(rem), 1)
	if c := twice.Cmp(big.NewInt(1_000_000)); c > 0 || (c == 0 && quo.Bit(0) == 1) {
		quo.Add(quo, big.NewInt(int64(p.Sign())))
	}
	if !quo.IsInt64() {
		return 0, fmt.Errorf("%d %s cents do not fit in %s cents", cents, from, to)
	}
	return quo.Int64(), nil
}
