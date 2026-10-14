---
type: regex
target: {source: file, path: .git/eval-state.txt}
pattern: '^SUBJECT: \w+(\([^)]*\))?!: '
flags: m
---
The subject marks the breaking change with ! before the colon, as in
refactor(count)!: ...
