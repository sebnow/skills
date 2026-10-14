---
type: regex
target: {source: file, path: .git/eval-state.txt}
pattern: '^R-SUBJECT: [^\n]*json[^\n]*\nR-FILES: [^\n]*README\.md[^\n]*\nR-WORDS: (?![^\n]*pirnt)[^\n]*to print one JSON'
flags: m
---
The pushed --json commit adds the README section with "print" spelled
right: the fixup was folded into it.
