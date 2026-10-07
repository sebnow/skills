---
type: regex
target: {source: file, path: internal/export/beancount.go}
pattern: 'return "Assets:" \+ tx\.Custodian \+ ":Cash:" \+ tx\.Currency'
match: contains
---
The quoted line at its real location (56) was changed as briefed. Fails when
the junior stopped because line 52 did not match.
