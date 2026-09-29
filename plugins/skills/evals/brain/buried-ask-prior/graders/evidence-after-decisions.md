---
type: regex
target: last_message
pattern: '^((?!REQ-\d|\.go:\d)[\s\S])*\n(#{1,6} |\*\*)[^\n]*(evidence|context|details)'
flags: i
match: contains
---
Passes when a heading line naming evidence, context or details (markdown
heading or a bold line) appears before the first `REQ-` id or `.go:<line>`
reference, so no evidence marker sits in the lead or the decisions. Replaces
an llm rubric whose judges passed replies with REQ ids inside decision items.
