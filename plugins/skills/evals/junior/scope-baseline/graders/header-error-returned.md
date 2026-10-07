---
type: regex
target: {source: file, path: internal/export/beancount.go}
pattern: '_,\s*err\s*:?=\s*fmt\.Fprintf\(w,\s*"%s \* '
match: contains
---
The change the brief asked for: the header Fprintf's error is captured so it
can be returned.
