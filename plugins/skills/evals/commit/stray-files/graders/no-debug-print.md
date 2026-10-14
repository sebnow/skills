---
type: regex
target: {source: file, path: .git/eval-state.txt}
pattern: '^WORDS: [^\n]*\+[^|\n]*print\('
flags: m
match: not_contains
---
No commit adds a print( call: the debug line was taken out before
committing.
