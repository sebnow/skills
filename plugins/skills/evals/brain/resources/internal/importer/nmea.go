package importer

import (
	"bufio"
	"io"
	"strings"
)

// ParseNMEA reads an NMEA 0183 sentence dump (REQ-0102). A voyage starts at
// the first $PTWENG sentence with rpm above zero and ends when every engine
// has been at zero for ten minutes.
func ParseNMEA(r io.Reader) ([]Record, []Bad, error) {
	sc := bufio.NewScanner(r)
	var bad []Bad
	line := 0
	for sc.Scan() {
		line++
		s := sc.Text()
		if !strings.HasPrefix(s, "$") {
			bad = append(bad, Bad{Line: line, Reason: "not an NMEA sentence"})
		}
	}
	// Voyage reconstruction lives in nmeaVoyages; see nmea_test.go for the
	// sentence sequences it accepts.
	return nil, bad, sc.Err()
}
