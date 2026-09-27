package fixture

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"tidewright/internal/flow"
	"tidewright/internal/importer"
	"tidewright/internal/ledger"
	"tidewright/internal/rules/ccvr"
)

func newLedger(t *testing.T) *ledger.Ledger {
	t.Helper()
	lg, err := ledger.Open(context.Background(), filepath.Join(t.TempDir(), "ledger.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { lg.Close() })
	return lg
}

func importAll(t *testing.T, lg *ledger.Ledger, f Fixture) {
	t.Helper()
	logs, err := filepath.Glob(filepath.Join(f.Path, "logs", "*"))
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range logs {
		if err := importer.ImportFile(context.Background(), lg, p); err != nil {
			t.Fatalf("import %s: %v", p, err)
		}
	}
}

// TestFixtures assesses each fixture's season from a fresh ledger.
func TestFixtures(t *testing.T) {
	fixtures, err := Load(filepath.Join("testdata", "fixtures"))
	if err != nil {
		t.Fatal(err)
	}
	if len(fixtures) == 0 {
		t.Fatal("no fixtures under testdata/fixtures")
	}
	for _, f := range fixtures {
		t.Run(filepath.Base(f.Path), func(t *testing.T) {
			lg := newLedger(t)
			importAll(t, lg, f)
			a, err := flow.AssessSeason(context.Background(), lg, f.Season)
			if err != nil {
				t.Fatal(err)
			}
			check(t, a, f.Expect)
		})
	}
}

// TestCoverage lists requirements and clauses no fixture covers. It fails only
// with FIXTURE_COVERAGE=strict, so a new requirement does not break the build.
func TestCoverage(t *testing.T) {
	fixtures, err := Load(filepath.Join("testdata", "fixtures"))
	if err != nil {
		t.Fatal(err)
	}
	covered := map[string]bool{}
	for _, f := range fixtures {
		for _, c := range append(f.Covers, f.Cites...) {
			covered[c] = true
		}
	}
	reqs, _ := filepath.Glob(filepath.Join("..", "..", "docs", "requirements", "REQ-*.md"))
	for _, r := range reqs {
		id := filepath.Base(r)[:len("REQ-0000")]
		if !covered[id] {
			if os.Getenv("FIXTURE_COVERAGE") == "strict" {
				t.Errorf("%s: no fixture covers it", id)
			} else {
				t.Logf("%s: no fixture covers it", id)
			}
		}
	}
}

// TestFixtureFiles checks each fixture has logs, a summary and a season the
// Regulations apply to, so a broken fixture fails here rather than in the
// rules.
func TestFixtureFiles(t *testing.T) {
	t.Parallel()
	fixtures, err := Load(filepath.Join("testdata", "fixtures"))
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range fixtures {
		name := filepath.Base(f.Path)
		logs, _ := filepath.Glob(filepath.Join(f.Path, "logs", "*"))
		if len(logs) == 0 {
			t.Errorf("%s: no files under logs/", name)
		}
		if f.Summary == "" {
			t.Errorf("%s: no summary", name)
		}
		if f.Season < 2022 {
			t.Errorf("%s: season %d is before CCVR 2021 came into force", name, f.Season)
		}
		if len(f.Expect) == 0 {
			t.Errorf("%s: expects nothing", name)
		}
	}
}

func check(t *testing.T, a *ccvr.Assessment, want []Expectation) {
	t.Helper()
	for _, e := range want {
		switch e.What {
		case "hours-towards-overhaul":
			h, ok := a.Engine(e.Engine)
			if !ok {
				t.Errorf("engine %s: not in the assessment", e.Engine)
			} else if h.SinceOverhaul != e.Hours {
				t.Errorf("engine %s: %.1f hours towards overhaul, want %.1f", e.Engine, h.SinceOverhaul, e.Hours)
			}
		case "survey-due":
			if got := a.SurveyDue.Format(time.DateOnly); got != e.Date {
				t.Errorf("survey due %s, want %s", got, e.Date)
			}
		case "levy-band":
			if int(a.Band) != e.Band {
				t.Errorf("levy band %d, want %d", a.Band, e.Band)
			}
		default:
			t.Errorf("unknown expectation %q", e.What)
		}
	}
}
