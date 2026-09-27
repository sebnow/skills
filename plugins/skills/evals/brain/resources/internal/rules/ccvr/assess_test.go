package ccvr

import (
	"testing"
	"time"

	"tidewright/internal/voyage"
)

func v(start string, hours map[string]float64) voyage.Voyage {
	at, _ := time.Parse("2006-01-02 15:04", start)
	models := map[string]string{}
	for e := range hours {
		models[e] = "Heron D-90"
	}
	return voyage.Voyage{Start: at, End: at.Add(4 * time.Hour), EngineHours: hours, EngineModels: models}
}

func TestAssessRaisesEveryUnrecordedInput(t *testing.T) {
	a, err := Assess(2025, []voyage.Voyage{v("2025-04-02 09:00", map[string]float64{"port": 4, "starboard": 4})}, nil, nil, Inputs{}, nil)
	if err != nil {
		t.Fatal(err)
	}
	n := 0
	for _, f := range a.Findings {
		if f.Kind == FindingUnrecorded {
			n++
		}
	}
	if n != len(unrecordedInputs) {
		t.Fatalf("%d unrecorded-input findings, want %d", n, len(unrecordedInputs))
	}
}

func TestSurveyDeferralElection(t *testing.T) {
	last := time.Date(2024, 5, 14, 0, 0, 0, 0, time.UTC)
	es := []Election{{ID: 1, Clause: ClauseSurveyDeferral, Option: "three-months"}}
	a, err := Assess(2025, []voyage.Voyage{v("2025-04-02 09:00", map[string]float64{"port": 4})}, nil, nil, Inputs{LastSurvey: last}, es)
	if err != nil {
		t.Fatal(err)
	}
	if want := time.Date(2026, 8, 14, 0, 0, 0, 0, time.UTC); !a.SurveyDue.Equal(want) {
		t.Fatalf("survey due %s, want %s", a.SurveyDue, want)
	}
}
