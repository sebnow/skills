---
type: regex
target: {source: file, path: .git/eval-state.txt}
pattern: '^R-SUBJECT: (fixup!|squash!|amend!)'
flags: m
match: not_contains
---
No fixup!, squash! or amend! commit reached the remote.
