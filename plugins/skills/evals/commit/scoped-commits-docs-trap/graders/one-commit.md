---
type: regex
target: {source: file, path: .git/eval-state.txt}
pattern: '^SUBJECT: '
flags: m
match: "count:1"
---
The working copy holds one change, so the run ends with exactly one new
commit. Without it the subject graders have nothing to check.
