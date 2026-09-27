package flow

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"tidewright/internal/ledger"
	"tidewright/internal/rules/ccvr"
	"tidewright/internal/voyage"
)

// Package flow joins the ledger to the rules: it reads what the ledger holds
// for a season, converts it to what ccvr.Assess takes, and runs the rules.
// Nothing it computes is written back (ADR-0004).

// firstSeason is the first season CCVR 2021 applies to: the Regulations came
// into force on 1 January 2022 (reg. 1(2)).
const firstSeason = 2022

// AssessSeason assesses one season from the ledger. The imported voyages are
// deduplicated first, so a clash is a finding rather than double-counted hours.
func AssessSeason(ctx context.Context, lg *ledger.Ledger, season int) (*ccvr.Assessment, error) {
	imports, err := lg.Imports(ctx)
	if err != nil {
		return nil, fmt.Errorf("flow: list imports: %w", err)
	}
	var voyages []voyage.Voyage
	var unread []ccvr.Unread
	for _, im := range imports {
		if im.Removed {
			continue
		}
		vs, bad, err := voyage.FromBlob(ctx, lg, im)
		if err != nil {
			return nil, fmt.Errorf("flow: read import %s: %w", im.Name, err)
		}
		voyages = append(voyages, vs...)
		for _, b := range bad {
			unread = append(unread, ccvr.Unread{Import: im.Name, Line: b.Line, Reason: b.Reason})
		}
	}
	voyages, clashes := voyage.Dedupe(voyages)

	overhauls, err := lg.Overhauls(ctx)
	if err != nil {
		return nil, fmt.Errorf("flow: overhauls: %w", err)
	}
	survey, err := lg.LastSurvey(ctx)
	if err != nil {
		return nil, fmt.Errorf("flow: last survey: %w", err)
	}
	inputs := ccvr.Inputs{LastSurvey: survey}
	for _, o := range overhauls {
		inputs.Overhauls = append(inputs.Overhauls, ccvr.Overhaul{Engine: o.Engine, At: o.At})
	}

	elections, err := electionsInForce(ctx, lg)
	if err != nil {
		return nil, err
	}
	if season < firstSeason {
		// A season before the Regulations came into force has no return to
		// make; refuse rather than produce figures nobody can submit.
		return nil, fmt.Errorf("flow: season %d is before CCVR 2021 came into force", season)
	}

	return ccvr.Assess(season, voyages, clashes, unread, inputs, elections)
}

// electionsInForce returns the option elections that no later election
// supersedes, converted for the rules.
func electionsInForce(ctx context.Context, lg *ledger.Ledger) ([]ccvr.Election, error) {
	all, err := lg.Elections(ctx)
	if err != nil {
		return nil, fmt.Errorf("flow: elections: %w", err)
	}
	var out []ccvr.Election
	for _, e := range all {
		// Only option elections reach the rules; superseded ones stay on record.
		if e.Kind != ledger.KindOption || e.SupersededBy != 0 {
			continue
		}
		var b ledger.OptionBody
		if err := json.Unmarshal(e.Body, &b); err != nil {
			return nil, fmt.Errorf("flow: election %d: %w", e.ID, err)
		}
		if b.Clause == "" {
			continue
		}
		out = append(out, ccvr.Election{ID: e.ID, Clause: b.Clause, Option: b.Option, MadeAt: e.MadeAt})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Clause < out[j].Clause })
	return out, nil
}
