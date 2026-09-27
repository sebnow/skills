package importer

import (
	"strings"
	"testing"
)

func TestParseNMEARejectsNonSentences(t *testing.T) {
	_, bad, err := ParseNMEA(strings.NewReader("$GPRMC,...\nnot a sentence\n"))
	if err != nil {
		t.Fatal(err)
	}
	if len(bad) != 1 || bad[0].Line != 2 {
		t.Fatalf("bad = %+v", bad)
	}
}
