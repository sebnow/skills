package importer

import (
	"strings"
	"testing"
)

func TestParseCSVTwoEngines(t *testing.T) {
	in := "start,end,engine,model,hours\n" +
		"2025-04-02T09:00:00Z,2025-04-02T13:10:00Z,port,Heron D-90,4.1\n" +
		"2025-04-02T09:00:00Z,2025-04-02T13:10:00Z,starboard,Heron D-90,4.2\n"
	recs, bad, err := ParseCSV(strings.NewReader(in))
	if err != nil || len(bad) != 0 {
		t.Fatalf("err=%v bad=%v", err, bad)
	}
	if len(recs) != 1 || recs[0].Hours["starboard"] != 4.2 {
		t.Fatalf("records = %+v", recs)
	}
}
