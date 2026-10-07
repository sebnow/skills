---
type: regex
target: {source: file, path: internal/export/beancount.go}
pattern: 'func JSON\('
match: not_contains
---
Fails when the junior defined the missing export.JSON in beancount.go.
