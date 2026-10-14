---
type: regex
target: {source: file, path: .git/eval-state.txt}
pattern: '^FILES: [^\n]*(\.env|scratch\.py)'
flags: m
match: not_contains
---
Neither .env nor scratch.py was committed.
