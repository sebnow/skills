// Package fetch gets the documents a return reads from outside the ledger.
// Each is fetched while the return is assessed and is not stored.
package fetch

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"tidewright/internal/vessel"
)

// StatusReport gets a classed vessel's latest status report from the address
// in its class entry. The overhaul and survey dates in it are used as given.
func StatusReport(ctx context.Context, v vessel.File) ([]byte, error) {
	if v.Class.ReportURL == "" {
		return nil, nil
	}
	return get(ctx, v.Class.ReportURL)
}

// BerthLet asks the harbour register the vessel file names whether the
// vessel's berth is let for the season. The answer decides the mooring link.
func BerthLet(ctx context.Context, v vessel.File, season int) ([]byte, error) {
	return get(ctx, fmt.Sprintf("%s/berths/%s/%d", v.Harbour.Register, v.Harbour.Berth, season))
}

// EngineHours gets an engine's hours from the maker's telemetry service at
// the address in the engine's record. When it answers, its figure replaces
// the hours summed from the imported logs.
func EngineHours(ctx context.Context, e vessel.Engine, season int) ([]byte, error) {
	if e.Telemetry == "" {
		return nil, nil
	}
	return get(ctx, fmt.Sprintf("%s/hours?season=%d", e.Telemetry, season))
}

// Amendments gets the amendments published after the regulation text in
// docs/regulation, from the feed the vessel file names. The rules apply them
// over the text in the repository.
func Amendments(ctx context.Context, v vessel.File) ([]byte, error) {
	return get(ctx, v.Regulation.Amendments+"/since/2021-04-01")
}

func get(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %s", url, resp.Status)
	}
	return io.ReadAll(resp.Body)
}
