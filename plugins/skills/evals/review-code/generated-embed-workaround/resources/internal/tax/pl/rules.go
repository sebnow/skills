package pl

import "fmt"

//go:generate go run ../../../cmd/gen-rulings

// Rule is a tax rule and the public ruling it relies on.
type Rule struct {
	ID        string
	RulingSig string
}

var Rules = []Rule{
	{ID: "income-on-settlement", RulingSig: "0112-KDIL2-1.4011.123.2025.1.AB"},
	{ID: "fifo-per-account", RulingSig: "DD8.8203.1.2021"},
}

// Evidence returns the text of the ruling r relies on, for inclusion in the
// evidence report.
func Evidence(r Rule) (string, error) {
	text, ok := rulingText[r.RulingSig]
	if !ok {
		return "", fmt.Errorf("rule %s: no text for ruling %s", r.ID, r.RulingSig)
	}
	return text, nil
}
