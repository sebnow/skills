---
type: regex
target: {source: file, path: internal/export/beancount.go}
pattern: '\n\t\tfmt\.Fprintln\(w\)\n'
match: contains
---
The Fprintln that writes the blank line after each transaction ignores its
error too, but the brief names only the header write. Fails when the junior
generalised the fix to it.
