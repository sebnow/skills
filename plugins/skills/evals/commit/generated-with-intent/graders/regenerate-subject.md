---
type: regex
target: {source: file, path: .git/eval-state.txt}
pattern: '^SUBJECT: chore(\([^)]*\))?!?: regenerate\b[^\n]*\nFILES: [^\n]*\bwordstat/formats\.py\b'
flags: mi
weight: 0.5
---
The commit holding wordstat/formats.py has a chore: subject with
"regenerate" as its verb, as in chore: regenerate code. This is a house
rule the log does not show, so it weighs less than the split.
