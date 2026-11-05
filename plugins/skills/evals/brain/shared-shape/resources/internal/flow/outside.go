package flow

import (
	"context"
	"fmt"
	"os"

	"tidewright/internal/fetch"
	"tidewright/internal/rules/ccvr"
	"tidewright/internal/vessel"
)

// fetchOutside reads the vessel file and fetches the documents it names into
// in.Outside. Nothing fetched is written to the ledger.
func fetchOutside(ctx context.Context, season int, in *ccvr.Inputs) error {
	path := os.Getenv("TIDEWRIGHT_VESSEL")
	if path == "" {
		path = "vessel.toml"
	}
	v, err := vessel.Load(path)
	if err != nil {
		return fmt.Errorf("flow: vessel file: %w", err)
	}
	if in.Outside.StatusReport, err = fetch.StatusReport(ctx, v); err != nil {
		return fmt.Errorf("flow: status report: %w", err)
	}
	if in.Outside.BerthLet, err = fetch.BerthLet(ctx, v, season); err != nil {
		return fmt.Errorf("flow: berth: %w", err)
	}
	in.Outside.EngineHours = map[string][]byte{}
	for _, e := range v.Engines {
		b, err := fetch.EngineHours(ctx, e, season)
		if err != nil {
			return fmt.Errorf("flow: engine %s hours: %w", e.ID, err)
		}
		if b != nil {
			in.Outside.EngineHours[e.ID] = b
		}
	}
	if in.Outside.Amendments, err = fetch.Amendments(ctx, v); err != nil {
		return fmt.Errorf("flow: amendments: %w", err)
	}
	return nil
}
