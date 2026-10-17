// Package rates holds the exchange rates the finance features convert with,
// loaded from the treasury team's daily export.
package rates

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

var ErrNoRate = errors.New("no rate for currency pair")

// Supported lists the currencies the export may quote. All of them have two
// decimal places.
var Supported = map[string]bool{"CHF": true, "EUR": true, "GBP": true, "PLN": true, "SEK": true, "USD": true}

// Quote is one published rate: one unit of From buys RateMicro/1,000,000
// units of To.
type Quote struct {
	From, To  string
	RateMicro int64
	AsOf      time.Time
}

// Book holds the latest quote for each currency pair. It does not change
// after Load.
type Book struct {
	quotes map[[2]string]Quote
}

// Load reads the treasury export: a header "as_of,from,to,rate", then one row
// per quote, with as_of in RFC 3339 and rate a positive decimal with at most
// six decimal places.
func Load(r io.Reader) (*Book, error) {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = 4
	header, err := cr.Read()
	if err != nil {
		return nil, fmt.Errorf("rates export: %w", err)
	}
	if strings.Join(header, ",") != "as_of,from,to,rate" {
		return nil, fmt.Errorf("rates export: unexpected header %q", strings.Join(header, ","))
	}
	b := &Book{quotes: map[[2]string]Quote{}}
	for line := 2; ; line++ {
		rec, err := cr.Read()
		if err == io.EOF {
			return b, nil
		}
		if err != nil {
			return nil, fmt.Errorf("rates export: %w", err)
		}
		q, err := parseQuote(rec)
		if err != nil {
			return nil, fmt.Errorf("rates export line %d: %w", line, err)
		}
		key := [2]string{q.From, q.To}
		if cur, ok := b.quotes[key]; !ok || q.AsOf.After(cur.AsOf) {
			b.quotes[key] = q
		}
	}
}

// Quote returns the latest quote for converting from one currency to another.
func (b *Book) Quote(from, to string) (Quote, error) {
	q, ok := b.quotes[[2]string{from, to}]
	if !ok {
		return Quote{}, fmt.Errorf("%w %s/%s", ErrNoRate, from, to)
	}
	return q, nil
}

func parseQuote(rec []string) (Quote, error) {
	asOf, err := time.Parse(time.RFC3339, rec[0])
	if err != nil {
		return Quote{}, fmt.Errorf("as_of: %w", err)
	}
	from, to := rec[1], rec[2]
	if !Supported[from] || !Supported[to] || from == to {
		return Quote{}, fmt.Errorf("pair %s/%s is not two different supported currencies", from, to)
	}
	rate, err := parseMicro(rec[3])
	if err != nil {
		return Quote{}, fmt.Errorf("rate %q: %w", rec[3], err)
	}
	return Quote{From: from, To: to, RateMicro: rate, AsOf: asOf.UTC()}, nil
}

// parseMicro parses a positive decimal with at most six decimal places into
// millionths.
func parseMicro(s string) (int64, error) {
	whole, frac, dot := strings.Cut(s, ".")
	switch {
	case !digits(whole) || (dot && !digits(frac)):
		return 0, errors.New("not a decimal number")
	case len(frac) > 6:
		return 0, errors.New("more than six decimal places")
	case len(whole) > 9:
		return 0, errors.New("too large")
	}
	n, err := strconv.ParseInt(whole+frac+strings.Repeat("0", 6-len(frac)), 10, 64)
	if err != nil {
		return 0, err
	}
	if n == 0 {
		return 0, errors.New("not positive")
	}
	return n, nil
}

func digits(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
