---
type: regex
target: {source: file, path: .jj/eval-state.txt}
pattern: '^FILES: README\.md tests/test_cli\.py wordstat/cli\.py$'
flags: m
---
One commit holds the --json flag with its test and README section, and
nothing else.
