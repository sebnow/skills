---
type: regex
target: {source: file, path: .git/eval-state.txt}
pattern: '^R-RESULT: (none|folded)$'
flags: m
---
Either nothing was pushed (asking first is fine), or the pushed feature
commit already holds the corrected README: the fixup was folded in, not
pushed or left for a later commit.
