package importer

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"time"
)

// Record is one voyage as an importer reads it.
type Record struct {
	Start, End time.Time
	Hours      map[string]float64
	Models     map[string]string
}

// Bad is a line the importer could not read (REQ-0126).
type Bad struct {
	Line   int
	Reason string
}

// ParseCSV reads a chartplotter CSV export (REQ-0101). Columns: start, end,
// engine, model, hours; one row per engine per voyage.
func ParseCSV(r io.Reader) ([]Record, []Bad, error) {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1
	rows, err := cr.ReadAll()
	if err != nil {
		return nil, nil, err
	}
	byStart := map[string]*Record{}
	var order []string
	var bad []Bad
	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) != 5 {
			bad = append(bad, Bad{Line: i + 1, Reason: fmt.Sprintf("%d columns, want 5", len(row))})
			continue
		}
		start, err1 := time.Parse(time.RFC3339, row[0])
		end, err2 := time.Parse(time.RFC3339, row[1])
		hours, err3 := strconv.ParseFloat(row[4], 64)
		if err1 != nil || err2 != nil || err3 != nil {
			bad = append(bad, Bad{Line: i + 1, Reason: "unreadable time or hours"})
			continue
		}
		rec, ok := byStart[row[0]]
		if !ok {
			rec = &Record{Start: start, End: end, Hours: map[string]float64{}, Models: map[string]string{}}
			byStart[row[0]] = rec
			order = append(order, row[0])
		}
		rec.Hours[row[2]] += hours
		rec.Models[row[2]] = row[3]
	}
	out := make([]Record, len(order))
	for i, k := range order {
		out[i] = *byStart[k]
	}
	return out, bad, nil
}
