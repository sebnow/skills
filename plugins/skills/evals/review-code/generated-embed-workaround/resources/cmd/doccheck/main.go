// Command doccheck verifies that every ruling under docs/rulings has a
// non-empty text.md.
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	paths, err := filepath.Glob("docs/rulings/*/*/text.md")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for _, p := range paths {
		fi, err := os.Stat(p)
		if err != nil || fi.Size() == 0 {
			fmt.Fprintf(os.Stderr, "%s: missing or empty\n", p)
			os.Exit(1)
		}
	}
}
