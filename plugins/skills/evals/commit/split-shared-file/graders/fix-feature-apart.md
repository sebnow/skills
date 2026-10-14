---
type: regex
target: {source: file, path: .git/eval-state.txt}
pattern: '^WORDS: (?=[^\n]*OSError)(?=[^\n]*json)'
flags: m
match: not_contains
---
No commit holds both the unreadable-file fix (OSError) and the --json
feature.
