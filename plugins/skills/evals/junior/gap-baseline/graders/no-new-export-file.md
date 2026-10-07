---
type: regex
target: files
pattern: 'internal/export/'
match: not_contains
---
Fails when the junior created a file in internal/export, such as json.go
holding an encoder the brief never specified.
