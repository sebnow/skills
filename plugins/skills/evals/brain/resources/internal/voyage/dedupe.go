package voyage

import (
	"sort"
	"time"
)

// Clash is two voyages from different imports that overlap in time
// (REQ-0108). The earlier-imported voyage is kept; the later one is dropped
// from the hours and reported.
type Clash struct {
	A, B     Voyage
	From, To time.Time
}

// Dedupe drops voyages that overlap an earlier one and returns the clashes.
func Dedupe(vs []Voyage) ([]Voyage, []Clash) {
	sort.SliceStable(vs, func(i, j int) bool { return vs[i].Start.Before(vs[j].Start) })
	var kept []Voyage
	var clashes []Clash
	for _, v := range vs {
		if n := len(kept); n > 0 && v.Start.Before(kept[n-1].End) && v.Import != kept[n-1].Import {
			to := v.End
			if kept[n-1].End.Before(to) {
				to = kept[n-1].End
			}
			clashes = append(clashes, Clash{A: kept[n-1], B: v, From: v.Start, To: to})
			continue
		}
		kept = append(kept, v)
	}
	return kept, clashes
}
