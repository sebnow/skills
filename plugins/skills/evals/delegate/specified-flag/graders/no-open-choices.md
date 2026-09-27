---
type: regex
target: trace
pattern: '"name":"Agent","input":\{.*?(your call|you decide|up to you|use your judge?ment)'
flags: i
match: not_contains
---
