package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"tidewright/internal/flow"
	"tidewright/internal/ledger"
	"tidewright/internal/rules/ccvr"
)

type command func(ctx context.Context, lg *ledger.Ledger, args []string, out io.Writer) error

var commands = map[string]command{
	"import":    runImport,
	"remove":    runRemove,
	"return":    runReturn,
	"status":    runStatus,
	"elect":     runElect,
	"elections": runElections,
	"clauses":   runClauses,
	"overhaul":  runOverhaul,
	"survey":    runSurvey,
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: tidewright <command> [args]; commands: import remove return status elect elections clauses overhaul survey")
		os.Exit(2)
	}
	cmd, ok := commands[os.Args[1]]
	if !ok {
		fmt.Fprintf(os.Stderr, "tidewright: unknown command %q\n", os.Args[1])
		os.Exit(2)
	}
	ctx := context.Background()
	path := os.Getenv("TIDEWRIGHT_LEDGER")
	if path == "" {
		path = "tidewright.db"
	}
	lg, err := ledger.Open(ctx, path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "tidewright:", err)
		os.Exit(1)
	}
	defer lg.Close()
	if err := cmd(ctx, lg, os.Args[2:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "tidewright:", err)
		os.Exit(1)
	}
}

func runReturn(ctx context.Context, lg *ledger.Ledger, args []string, out io.Writer) error {
	season, err := seasonArg(args)
	if err != nil {
		return err
	}
	a, err := flow.AssessSeason(ctx, lg, season)
	if err != nil {
		return err
	}
	return ccvr.RenderReturn(out, a)
}
