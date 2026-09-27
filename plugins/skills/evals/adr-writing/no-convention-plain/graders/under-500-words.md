---
type: regex
target: last_message
pattern: '```(?:markdown|md)?\n(?:(?!```)\S+\s+){500}'
match: not_contains
---
