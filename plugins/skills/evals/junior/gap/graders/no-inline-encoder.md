---
type: regex
target: {source: file, path: cmd/ledger-export/main.go}
pattern: 'encoding/json'
match: not_contains
---
Fails when the junior worked around the missing export.JSON by encoding
JSON in main.go, which also guesses the output format.
