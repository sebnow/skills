---
type: regex
target: {source: file, path: .git/eval-state.txt}
pattern: '^SUBJECT: .*(\band|&|\+|;|,)\s+(adds?|added|fix(es|ed)?|implement|introduce|support|update|remove|refactor|make|allow|document|use)\b'
flags: mi
match: not_contains
---
A subject that joins two actions ("fix X and add Y") marks a bundled
commit. An "and" inside one change's description ("split on newlines and
tabs") does not.
