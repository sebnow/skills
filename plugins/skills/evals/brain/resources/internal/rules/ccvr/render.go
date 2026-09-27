package ccvr

import (
	"fmt"
	"io"
	"strings"
)

// The return is Markdown: a status banner, the figures, the elections in
// force, and the open findings, each explained as ADR-0006 requires.

const disclaimer = "Prepared by tidewright from the imported voyage logs. Check every figure against the logs before you submit."

// section writes a level-two heading followed by a blank line.
func section(w io.Writer, title string) {
	fmt.Fprintf(w, "\n## %s\n\n", title)
}

// cell escapes a value for a Markdown table cell.
func cell(s string) string {
	return strings.ReplaceAll(s, "|", `\|`)
}

// hours formats engine hours the way the operating return prints them:
// one decimal place.
func hours(h float64) string {
	return fmt.Sprintf("%.1f", h)
}

// RenderReturn writes the operating return for an assessment. The first line
// is the status banner: a return with open findings says so before anything
// else (REQ-0137).
func RenderReturn(w io.Writer, a *Assessment) error {
	p := func(format string, args ...any) { fmt.Fprintf(w, format+"\n", args...) }
	if a.Submittable() {
		p("No open findings. The figures may be entered on the CCVR operating return.")
	} else {
		p("**NOT FOR SUBMISSION**: %d open findings, listed under \"Open findings\" below.", len(a.Findings))
	}
	p("")
	p("# CCVR operating return, season %d", a.Season)
	p("")
	p("_%s_", disclaimer)

	section(w, "Engine hours")
	p("| Engine | Class | Hours since overhaul | Threshold |")
	p("|---|---|---|---|")
	for _, h := range a.Engines {
		p("| %s | %s | %s | %s |", cell(h.Engine), h.Class, hours(h.SinceOverhaul), hours(h.Threshold))
	}
	p("")
	p("Engine hours in the season: %s", hours(a.SeasonHours))

	section(w, "Survey")
	p("Next survey due: %s", a.SurveyDue.Format("2 January 2006"))

	section(w, "Levy")
	p("Levy band: %d (Schedule 3)", a.Band)

	if len(a.Elections) > 0 {
		section(w, "Elections in force")
		p("| Clause | Option | Made |")
		p("|---|---|---|")
		for _, e := range a.Elections {
			p("| %s | %s | %s |", cell(e.Clause), cell(e.Option), e.MadeAt.Format("2006-01-02"))
		}
	}

	renderEngineNotes(w, a)
	if !a.Submittable() {
		renderFindings(w, a, p)
		renderCitations(w, a)
	}
	return nil
}

var monthNames = [12]string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}

// RenderMonthly writes the season's engine hours by month. The surveyor's
// form asks for it on the second page.
func RenderMonthly(w io.Writer, vs []VoyageHours) {
	p := func(format string, args ...any) { fmt.Fprintf(w, format+"\n", args...) }
	m := monthHours(vs)
	section(w, "Engine hours by month")
	p("| Month | Hours |")
	p("|---|---|")
	for i, h := range m {
		if h == 0 {
			continue
		}
		p("| %s | %s |", monthNames[i], hours(h))
	}
}

// RenderSummary writes the one-line status `tidewright status` prints.
func RenderSummary(w io.Writer, a *Assessment) {
	if a.Submittable() {
		fmt.Fprintf(w, "season %d: submittable\n", a.Season)
		return
	}
	kinds := map[FindingKind]int{}
	for _, f := range a.Findings {
		kinds[f.Kind]++
	}
	var parts []string
	for _, k := range []FindingKind{FindingClash, FindingUnread, FindingOverhaul, FindingElection, FindingUnrecorded} {
		if kinds[k] > 0 {
			parts = append(parts, fmt.Sprintf("%d %s", kinds[k], k))
		}
	}
	fmt.Fprintf(w, "season %d: NOT FOR SUBMISSION (%s)\n", a.Season, strings.Join(parts, ", "))
}

// RenderFindingsText writes the open findings as plain text, one per
// paragraph, for `tidewright findings`.
func RenderFindingsText(w io.Writer, a *Assessment) {
	for i, f := range a.Findings {
		if i > 0 {
			fmt.Fprintln(w)
		}
		fmt.Fprintf(w, "%s\n  %s\n", f.Label, f.Text)
		if f.Assumes != "" {
			fmt.Fprintf(w, "  %s\n", f.Assumes)
		}
		fmt.Fprintf(w, "  fix: %s\n", f.Fix)
	}
}

// renderCitations lists every clause the return relies on, once each, so the
// surveyor can check them against the Regulations.
func renderCitations(w io.Writer, a *Assessment) {
	seen := map[string]bool{}
	var list []string
	for _, f := range a.Findings {
		for _, c := range f.Cites {
			if !seen[c.Clause] {
				seen[c.Clause] = true
				list = append(list, c.Clause)
			}
		}
	}
	if len(list) == 0 {
		return
	}
	section(w, "Clauses relied on")
	for _, c := range list {
		fmt.Fprintf(w, "- %s\n", c)
	}
}

// renderEngineNotes adds a line for each engine within 10 per cent of its
// threshold, so the operator can book the overhaul before it is overdue.
func renderEngineNotes(w io.Writer, a *Assessment) {
	for _, h := range a.Engines {
		if h.Over() || h.SinceOverhaul < 0.9*h.Threshold {
			continue
		}
		fmt.Fprintf(w, "\nEngine %s is within %s hours of its overhaul threshold.\n", h.Engine, hours(h.Threshold-h.SinceOverhaul))
	}
}

// renderAssumptions lists what the return assumed in place of missing inputs.
func renderAssumptions(w io.Writer, a *Assessment) {
	fmt.Fprintln(w, "Assumed in place of missing inputs:")
	for _, f := range a.Findings {
		if f.Assumes != "" {
			fmt.Fprintf(w, "- %s\n", f.Assumes)
		}
	}
}

// renderFindings writes the open findings, grouped by kind, each with what it
// assumes, the clause it rests on, and how to fix it (ADR-0006).
func renderFindings(w io.Writer, a *Assessment, p func(string, ...any)) {
	p("")
	p("## Open findings")
	p("")
	p("Each finding below blocks submission until it is resolved.")
	p("")
	n := 0
	for _, f := range a.Findings {
		n++
		p("%d. **%s.** %s", n, f.Label, f.Text)
		if f.Assumes != "" {
			p("   - Assumed: %s", f.Assumes)
		}
		for _, c := range f.Cites {
			if c.Text != "" {
				p("   - %s: \"%s\"", c.Clause, c.Text)
			} else {
				p("   - %s", c.Clause)
			}
		}
		p("   - Fix: %s", f.Fix)
	}
}
