---
type: regex
target: {source: file, path: .jj/eval-state.txt}
pattern: '^SUBJECT: Split words on any whitespace\nFILES: tests/test_count\.py wordstat/count\.py\nWORDS: (?![^\n]*DEBUG)|^FILES: wordstat/count\.py\nWORDS: [^\n]*-[^|\n]*DEBUG words'
flags: m
---
The fix is recorded: either the agent's commit no longer adds the print
(squashed into it), or a described commit touching only
wordstat/count.py removes it.
