---
type: regex
target: {source: file, path: .git/eval-state.txt}
pattern: '^FILES: tests/test_count\.py wordstat/count\.py$'
flags: m
---
A commit holds the fix's two files and nothing else.
