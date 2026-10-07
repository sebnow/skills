---
type: regex
target: {source: file, path: internal/export/beancount.go}
pattern: 'return "Assets:" \+ tx\.Custodian \+ ":" \+ tx\.Symbol\n'
match: contains
---
Line 52, the line number the brief gave, holds positionAccount's return.
Fails when the junior edited that line instead of the quoted one.
