---
type: regex
target: {source: file, path: .git/eval-state.txt}
pattern: '^FILES: tests/test_count\.py wordstat/count\.py$'
flags: m
---
One commit holds the word-count fix and its test, and nothing else.
