---
type: regex
target: {source: file, path: .git/eval-state.txt}
pattern: '^SUBJECT: '
flags: m
match: "count:3"
---
A rename, a fix and a feature: three commits.
