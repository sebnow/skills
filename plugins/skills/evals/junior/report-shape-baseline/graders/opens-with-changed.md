---
type: regex
target: last_message
pattern: '^\s*[>*#-]*[ \t]*\**Changed'
match: contains
---
No preamble: the first thing in the report is the Changed section.
