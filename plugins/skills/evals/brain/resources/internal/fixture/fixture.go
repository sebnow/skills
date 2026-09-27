package fixture

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// A fixture is a season of real-shaped voyage logs with the figures the
// return must produce for it. Each fixture lives in
// testdata/fixtures/<name>/ with its log files and a fixture.yaml.

// Fixture describes one fixture. Covers and Cites say what it is evidence for,
// so a requirement or clause with no fixture shows up in `go test -run Coverage`.
type Fixture struct {
	// Path is the fixture directory, set by Load.
	Path string `yaml:"-"`
	// Summary says in one line what the season contains.
	Summary string `yaml:"summary"`
	Season  int    `yaml:"season"`
	// Covers lists the REQ ids the fixture is evidence for.
	Covers []string `yaml:"covers"`
	// Cites lists the CCVR clauses the fixture exercises.
	Cites []string `yaml:"cites"`
	// Expect lists the figures the return must contain.
	Expect []Expectation `yaml:"expect"`
}

// Expectation is one figure the return must contain.
type Expectation struct {
	What   string  `yaml:"what"`
	Engine string  `yaml:"engine,omitempty"`
	Hours  float64 `yaml:"hours,omitempty"`
	Date   string  `yaml:"date,omitempty"`
	Band   int     `yaml:"band,omitempty"`
}

// Load reads every fixture under dir.
func Load(dir string) ([]Fixture, error) {
	paths, err := filepath.Glob(filepath.Join(dir, "*", "fixture.yaml"))
	if err != nil {
		return nil, err
	}
	var out []Fixture
	for _, p := range paths {
		b, err := os.ReadFile(p)
		if err != nil {
			return nil, err
		}
		var f Fixture
		if err := yaml.Unmarshal(b, &f); err != nil {
			return nil, fmt.Errorf("fixture %s: %w", p, err)
		}
		f.Path = filepath.Dir(p)
		out = append(out, f)
	}
	return out, nil
}
