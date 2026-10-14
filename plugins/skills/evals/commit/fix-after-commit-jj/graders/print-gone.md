---
type: regex
target: {source: file, path: wordstat/count.py}
pattern: 'DEBUG'
match: not_contains
---
The debug print is gone from wordstat/count.py.
