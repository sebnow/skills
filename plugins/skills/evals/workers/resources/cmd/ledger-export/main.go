package main

import (
	"flag"
	"fmt"
	"os"

	"example.com/ledgerx/internal/export"
)

func main() {
	db := flag.String("db", "portfolio.sqlite", "portfolio database")
	out := flag.String("out", "ledger.beancount", "output file")
	flag.Parse()

	txs, err := export.Load(*db)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	f, err := os.Create(*out)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer f.Close()
	if err := export.Beancount(txs, f); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
