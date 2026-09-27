package voyage

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"tidewright/internal/importer"
	"tidewright/internal/ledger"
)

// Voyage is one passage from casting off to making fast.
type Voyage struct {
	Import       string
	Start, End   time.Time
	EngineHours  map[string]float64
	EngineModels map[string]string
}

// Bad is a record an importer could not read.
type Bad struct {
	Line   int
	Reason string
}

// FromBlob reads the voyages of one import.
func FromBlob(ctx context.Context, lg *ledger.Ledger, im ledger.Import) ([]Voyage, []Bad, error) {
	content, err := lg.Content(ctx, im.Blob)
	if err != nil {
		return nil, nil, err
	}
	var recs []importer.Record
	var bad []importer.Bad
	switch im.Format {
	case "csv":
		recs, bad, err = importer.ParseCSV(bytes.NewReader(content))
	case "nmea":
		recs, bad, err = importer.ParseNMEA(bytes.NewReader(content))
	default:
		return nil, nil, fmt.Errorf("voyage: import %s has unknown format %q", im.Name, im.Format)
	}
	if err != nil {
		return nil, nil, err
	}
	out := make([]Voyage, len(recs))
	for i, r := range recs {
		out[i] = Voyage{Import: im.Name, Start: r.Start, End: r.End, EngineHours: r.Hours, EngineModels: r.Models}
	}
	outBad := make([]Bad, len(bad))
	for i, b := range bad {
		outBad[i] = Bad{Line: b.Line, Reason: b.Reason}
	}
	return out, outBad, nil
}
