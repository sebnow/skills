package ccvr

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"tidewright/internal/voyage"
)

// FindingKind classifies a finding for the return's layout.
type FindingKind string

const (
	FindingClash      FindingKind = "clash"
	FindingUnread     FindingKind = "unread"
	FindingOverhaul   FindingKind = "overhaul"
	FindingUnrecorded FindingKind = "unrecorded"
	FindingElection   FindingKind = "election"
)

// Finding is something the return has to tell the reader. Every finding in
// an assessment is open: the return cannot be submitted while one remains
// (REQ-0137). ADR-0006 fixes what a finding says.
type Finding struct {
	Kind    FindingKind
	Label   string
	Text    string
	Assumes string
	Cites   []Citation
	Fix     string
}

// Assessment is what the rules compute for one season; nothing in it is stored.
type Assessment struct {
	Season      int
	Engines     []EngineHours
	SeasonHours float64
	SurveyDue   time.Time
	Band        LevyBand
	Elections   []Election
	Findings    []Finding
}

// VoyageHours is the part of a voyage the rules read.
type VoyageHours struct {
	Start  time.Time
	Hours  map[string]float64
	Models map[string]string
}

// Overhaul is a recorded engine overhaul.
type Overhaul struct {
	Engine string
	At     time.Time
}

// Unread is a record an importer could not read (REQ-0126).
type Unread struct {
	Import string
	Line   int
	Reason string
}

// Inputs is what the ledger holds about the vessel apart from voyages and
// elections.
type Inputs struct {
	LastSurvey time.Time
	Overhauls  []Overhaul
}

// ErrNoVoyages is returned when a season has no imported voyages.
var ErrNoVoyages = errors.New("ccvr: no voyages in season")

// Engine returns the hours towards overhaul for one engine.
func (a *Assessment) Engine(name string) (EngineHours, bool) {
	for _, h := range a.Engines {
		if h.Engine == name {
			return h, true
		}
	}
	return EngineHours{}, false
}

// Election is an election in force, as the rules see it. flow builds these
// from the ledger's KindOption elections.
type Election struct {
	ID     int64
	Clause string
	Option string
	MadeAt time.Time
}

// electionFor returns the election in force for clause, if any.
func electionFor(es []Election, clause string) (Election, bool) {
	for _, e := range es {
		if e.Clause == clause {
			return e, true
		}
	}
	return Election{}, false
}

// clashFinding explains two overlapping voyages (REQ-0108).
func clashFinding(c voyage.Clash) Finding {
	return Finding{
		Kind:  FindingClash,
		Label: "Overlapping voyages",
		Text: fmt.Sprintf("Open: voyages from %s and %s overlap between %s and %s; neither is counted twice",
			c.A.Import, c.B.Import, c.From.Format("2006-01-02 15:04"), c.To.Format("2006-01-02 15:04")),
		Fix: "Remove the wrong import with `tidewright remove`, then import the right file.",
	}
}

// unreadFinding explains a record an importer skipped (REQ-0126).
func unreadFinding(u Unread) Finding {
	return Finding{
		Kind:  FindingUnread,
		Label: "Unread record",
		Text:  fmt.Sprintf("Open: %s line %d could not be read (%s)", u.Import, u.Line, u.Reason),
		Fix:   "Fix the export on the chartplotter and import it again.",
	}
}

// overhaulFinding explains an engine past its Schedule 2 threshold.
func overhaulFinding(h EngineHours) Finding {
	return Finding{
		Kind:  FindingOverhaul,
		Label: "Overhaul overdue",
		Text: fmt.Sprintf("Open: engine %s has run %.1f hours since its last overhaul; the class %s threshold is %.0f (CCVR 2021 reg. 11(1))",
			h.Engine, h.SinceOverhaul, h.Class, h.Threshold),
		Cites: []Citation{{Clause: "CCVR 2021 reg. 11(1)", Text: quotes["CCVR 2021 reg. 11(1)"]}},
		Fix:   "Record the overhaul with `tidewright overhaul`.",
	}
}

// seasonVoyages converts the season's voyages for the rules.
func seasonVoyages(season int, vs []voyage.Voyage) []VoyageHours {
	out := make([]VoyageHours, 0, len(vs))
	for _, v := range vs {
		if v.Start.Year() != season {
			continue
		}
		out = append(out, VoyageHours{Start: v.Start, Hours: v.EngineHours, Models: v.EngineModels})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Start.Before(out[j].Start) })
	return out
}

// laidUpFinding explains a season with voyages that the operator elected to
// treat as laid up.
func laidUpFinding(season int) Finding {
	return Finding{
		Kind:  FindingElection,
		Label: "Laid-up season",
		Text:  fmt.Sprintf("Open: season %d is elected as laid up, but it has voyages", season),
		Cites: []Citation{{Clause: "CCVR 2021 reg. 19(1)", Text: quotes["CCVR 2021 reg. 19(1)"]}},
		Fix:   "Elect `laid-up-season counts`, or remove the imports for this season.",
	}
}

// unknownClauseFinding explains an election on a clause the rules no longer
// offer, for example after a clause was withdrawn from the registry.
func unknownClauseFinding(e Election) Finding {
	return Finding{
		Kind:  FindingElection,
		Label: "Election on an unknown clause",
		Text:  fmt.Sprintf("Open: election %d is on clause %q, which offers no election", e.ID, e.Clause),
		Fix:   "Check `tidewright elections` and elect again on a clause from `tidewright clauses`.",
	}
}

// findingOrder is the order findings appear in the return: problems with the
// logs first, then the rules, then what nothing records yet.
var findingOrder = map[FindingKind]int{
	FindingClash:      0,
	FindingUnread:     1,
	FindingOverhaul:   2,
	FindingElection:   3,
	FindingUnrecorded: 4,
}

// sortFindings orders findings for the return, keeping the order within a kind.
func sortFindings(fs []Finding) {
	sort.SliceStable(fs, func(i, j int) bool { return findingOrder[fs[i].Kind] < findingOrder[fs[j].Kind] })
}

// monthHours sums the season's engine hours by month, for the monthly table
// in the return. Hours of a voyage that crosses midnight at the end of a
// month count in the month it started.
func monthHours(vs []VoyageHours) [12]float64 {
	var out [12]float64
	for _, v := range vs {
		for _, h := range v.Hours {
			out[v.Start.Month()-1] += h
		}
	}
	return out
}

// Submittable reports whether the return can be submitted (REQ-0137).
func (a *Assessment) Submittable() bool {
	return len(a.Findings) == 0
}

// Assess applies the rules to one season. It never fails on a gap in the
// inputs: a gap becomes an open finding that names what the return assumed.
func Assess(season int, voyages []voyage.Voyage, clashes []voyage.Clash, unread []Unread, inputs Inputs, elections []Election) (*Assessment, error) {
	vs := seasonVoyages(season, voyages)
	if len(vs) == 0 {
		return nil, ErrNoVoyages
	}
	a := &Assessment{Season: season, Elections: elections}

	engines, err := hoursTowardsOverhaul(vs, inputs.Overhauls)
	if err != nil {
		return nil, err
	}
	a.Engines = engines
	for _, v := range vs {
		for _, h := range v.Hours {
			a.SeasonHours += h
		}
	}
	a.Band = levyBand(a.SeasonHours)

	a.SurveyDue = surveyDue(inputs.LastSurvey)
	if e, ok := electionFor(elections, ClauseSurveyDeferral); ok && e.Option == "three-months" {
		a.SurveyDue = a.SurveyDue.AddDate(0, 3, 0)
	}

	for _, c := range clashes {
		a.Findings = append(a.Findings, clashFinding(c))
	}
	for _, u := range unread {
		a.Findings = append(a.Findings, unreadFinding(u))
	}
	for _, h := range engines {
		if h.Over() {
			a.Findings = append(a.Findings, overhaulFinding(h))
		}
	}
	for _, e := range elections {
		if _, ok := ClauseByKey(e.Clause); !ok {
			a.Findings = append(a.Findings, unknownClauseFinding(e))
		}
	}
	if e, ok := electionFor(elections, ClauseLaidUpSeason); ok && e.Option == "excluded" {
		a.Findings = append(a.Findings, laidUpFinding(season))
	}
	if e, ok := electionFor(elections, ClauseOverhaulRecord); ok && e.Option == "per-vessel" {
		// reg. 12(2): a per-vessel record resets every engine at the
		// latest overhaul of any of them.
		var latest time.Time
		for _, o := range inputs.Overhauls {
			if o.At.After(latest) {
				latest = o.At
			}
		}
		var reset []Overhaul
		for _, h := range engines {
			reset = append(reset, Overhaul{Engine: h.Engine, At: latest})
		}
		if engines, err = hoursTowardsOverhaul(vs, reset); err != nil {
			return nil, err
		}
		a.Engines = engines
	}
	if e, ok := electionFor(elections, ClauseTenderAsVessel); ok && e.Option == "part-of-parent" {
		// reg. 20(2): the tender's outboard counts towards the parent's
		// season hours but not towards any engine's overhaul.
		for _, v := range vs {
			if h, ok := v.Hours["tender"]; ok {
				a.SeasonHours += h
			}
		}
		a.Band = levyBand(a.SeasonHours)
	}
	if len(a.Elections) > 0 {
		sort.Slice(a.Elections, func(i, j int) bool { return a.Elections[i].Clause < a.Elections[j].Clause })
	}

	for _, u := range unrecordedInputs {
		a.Findings = append(a.Findings, Finding{
			Kind:    FindingUnrecorded,
			Label:   u.Label,
			Text:    "Open: the return needs " + u.Needs + ", and nothing records it (" + clauses(u.Cites) + ")",
			Assumes: u.Assumes,
			Cites:   u.Cites,
			Fix:     "No command records this input yet.",
		})
	}

	return a, nil
}
