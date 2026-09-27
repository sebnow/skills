---
type: regex
target: last_message
pattern: '^\s*(#{1,4}\s*|\*\*)(Status|Context|Decision|Consequences)\b'
flags: im
match: not_contains
---
