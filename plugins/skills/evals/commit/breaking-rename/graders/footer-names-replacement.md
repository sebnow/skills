---
type: regex
target: {source: file, path: .git/eval-state.txt}
pattern: '^BODY: BREAKING[ -]CHANGE:(?:[^\n]*\nBODY: (?=[^\n]))*[^\n]*count_text'
flags: m
---
The BREAKING CHANGE footer, on its line or the lines that continue it
before a blank line, names count_text, so callers learn what to use
instead.
