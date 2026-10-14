---
type: regex
target: {source: file, path: .git/eval-state.txt}
pattern: '^BODY: BREAKING[ -]CHANGE: '
flags: m
---
The body has a footer line starting BREAKING CHANGE: or BREAKING-CHANGE:.
