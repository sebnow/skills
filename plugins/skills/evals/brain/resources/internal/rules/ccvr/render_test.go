package ccvr

import (
	"strings"
	"testing"
	"time"
)

func assessment(findings ...Finding) *Assessment {
	return &Assessment{
		Season:      2025,
		Engines:     []EngineHours{{Engine: "port", Class: EngineClassC, SinceOverhaul: 412.5, Threshold: 6000}},
		SeasonHours: 412.5,
		SurveyDue:   time.Date(2026, 5, 14, 0, 0, 0, 0, time.UTC),
		Band:        2,
		Findings:    findings,
	}
}

func TestRenderReturnStatus(t *testing.T) {
	open := Finding{Kind: FindingUnrecorded, Label: "Home port", Text: "Open: ...", Fix: "No command records this input yet."}
	tests := []struct {
		name        string
		a           *Assessment
		wantBlocked bool
	}{
		{"no findings", assessment(), false},
		{"one open finding", assessment(open), true},
		{"two open findings", assessment(open, open), true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var b strings.Builder
			if err := RenderReturn(&b, tc.a); err != nil {
				t.Fatal(err)
			}
			got := b.String()
			first, _, _ := strings.Cut(got, "\n")
			if tc.wantBlocked && !strings.HasPrefix(first, "**NOT FOR SUBMISSION**") {
				t.Errorf("first line %q, want the status banner", first)
			}
			if !strings.Contains(got, "Levy band: 2") {
				t.Errorf("levy band missing from:\n%s", got)
			}
			if !strings.Contains(got, "| port | C | 412.5 | 6000.0 |") {
				t.Errorf("engine table missing from:\n%s", got)
			}
			if tc.wantBlocked != strings.Contains(got, "## Open findings") {
				t.Errorf("open findings section present = %v, want %v", !tc.wantBlocked, tc.wantBlocked)
			}

			if strings.Contains(got, "NOT FOR SUBMISSION") != tc.wantBlocked {
				t.Errorf("blocked = %v, want %v", !tc.wantBlocked, tc.wantBlocked)
			}
		})
	}
}
