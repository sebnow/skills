---
type: regex
target: {source: file, path: .git/eval-state.txt}
pattern: '^FILES: (?![^\n]*\bdata/formats\.csv\b)[^\n]*\bwordstat/formats\.py\b'
flags: m
---
The commit holding the regenerated wordstat/formats.py does not also hold
data/formats.csv, the hand-written source it was generated from. Fails
when formats.py was not committed at all.
