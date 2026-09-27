package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"text/tabwriter"

	"tidewright/internal/ledger"
	"tidewright/internal/rules/ccvr"
)

// runElect stores which option the operator picks for a clause:
// tidewright elect <clause> <option> [-evidence FILE].
// A later election for the same clause supersedes the earlier one;
// neither is deleted (REQ-0139).
func runElect(ctx context.Context, lg *ledger.Ledger, args []string, out io.Writer) error {
	fs := flag.NewFlagSet("elect", flag.ContinueOnError)
	evidence := fs.String("evidence", "", "stored file supporting the election")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 2 {
		return errors.New("usage: tidewright elect <clause> <option> [-evidence FILE]")
	}
	id, option := fs.Arg(0), fs.Arg(1)
	cl, ok := ccvr.ClauseByKey(id)
	if !ok {
		return fmt.Errorf("no clause %q offers an election; see tidewright clauses", id)
	}
	if !cl.Offers(option) {
		return fmt.Errorf("clause %q has no option %q (options: %v)", id, option, cl.OptionKeys())
	}
	var sources []ledger.Source
	if *evidence != "" {
		blob, err := lg.BlobByName(ctx, *evidence)
		if err != nil {
			return fmt.Errorf("evidence %q: %w", *evidence, err)
		}
		sources = append(sources, ledger.Source{Blob: blob.ID})
	}
	e, err := lg.Elect(ctx, id, option, sources)
	if err != nil {
		return err
	}
	fmt.Fprintf(out, "election %d: %s = %s\n", e.ID, id, option)
	if e.Supersedes != 0 {
		fmt.Fprintf(out, "supersedes election %d\n", e.Supersedes)
	}
	return nil
}

// runElections lists elections, in force first, then superseded.
func runElections(ctx context.Context, lg *ledger.Ledger, args []string, out io.Writer) error {
	all, err := lg.Elections(ctx)
	if err != nil {
		return err
	}
	tw := tabwriter.NewWriter(out, 0, 4, 2, ' ', 0)
	fmt.Fprintln(tw, "ID\tKIND\tABOUT\tVALUE\tMADE\tSTATUS")
	for _, e := range all {
		status := "in force"
		if e.SupersededBy != 0 {
			status = fmt.Sprintf("superseded by %d", e.SupersededBy)
		}
		fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%s\t%s\n", e.ID, e.Kind, e.About, e.Summary(), e.MadeAt.Format("2006-01-02"), status)
	}
	return tw.Flush()
}

// runClauses lists the clauses that offer an election and their options.
func runClauses(ctx context.Context, lg *ledger.Ledger, args []string, out io.Writer) error {
	tw := tabwriter.NewWriter(out, 0, 4, 2, ' ', 0)
	fmt.Fprintln(tw, "CLAUSE\tLABEL\tOPTIONS")
	for _, c := range ccvr.Clauses() {
		fmt.Fprintf(tw, "%s\t%s\t%v\n", c.Key, c.Label, c.OptionKeys())
	}
	return tw.Flush()
}
