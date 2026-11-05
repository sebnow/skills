package ccvr

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// statusReport is the part of a classification society's status report the
// rules read.
type statusReport struct {
	SurveyDue   time.Time            `json:"survey_due"`
	OverhaulDue map[string]time.Time `json:"overhaul_due"`
}

// applyOutside applies the documents fetched during the assessment over the
// figures computed from the ledger. A fetched figure wins over a computed one.
func applyOutside(a *Assessment, o Outside) error {
	if o.StatusReport != nil {
		var r statusReport
		if err := json.Unmarshal(o.StatusReport, &r); err != nil {
			return fmt.Errorf("ccvr: status report: %w", err)
		}
		a.SurveyDue = r.SurveyDue
		a.OverhaulDue = r.OverhaulDue
	}
	if o.BerthLet != nil {
		let := strings.TrimSpace(string(o.BerthLet)) == "let"
		a.BerthLet = &let
	}
	for engine, b := range o.EngineHours {
		hours, err := strconv.ParseFloat(strings.TrimSpace(string(b)), 64)
		if err != nil {
			return fmt.Errorf("ccvr: telemetry hours for %s: %w", engine, err)
		}
		for i := range a.Engines {
			if a.Engines[i].Engine == engine {
				a.Engines[i].SinceOverhaul = hours
			}
		}
	}
	if o.Amendments != nil {
		var ams []struct {
			Ref string `json:"ref"`
		}
		if err := json.Unmarshal(o.Amendments, &ams); err != nil {
			return fmt.Errorf("ccvr: amendments: %w", err)
		}
		for _, am := range ams {
			a.Amendments = append(a.Amendments, am.Ref)
		}
	}
	return nil
}
