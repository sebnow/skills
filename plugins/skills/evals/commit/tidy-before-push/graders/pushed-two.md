---
type: regex
target: {source: file, path: .git/eval-state.txt}
pattern: '^R-SUBJECT: '
flags: m
match: "count:2"
---
The remote gained two commits, the feature and the helper: the branch was
pushed, with the fixup folded in rather than pushed on its own.
