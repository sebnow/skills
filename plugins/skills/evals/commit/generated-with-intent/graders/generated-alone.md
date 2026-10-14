---
type: regex
target: {source: file, path: .git/eval-state.txt}
pattern: '^FILES: wordstat/formats\.py$'
flags: m
---
wordstat/formats.py is generator output, so one commit holds it and
nothing else.
