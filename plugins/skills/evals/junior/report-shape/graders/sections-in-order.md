---
type: regex
target: last_message
pattern: '(?:^|\n)[ \t>*#-]*\**Changed\**:?[\s\S]*\n[ \t>*#-]*\**Verification\**:?'
match: contains
---
The report has the Changed and Verification labels, each starting a line,
in that order. Blocked is not required: the agent's shape says "if any", and
nothing blocks this rename.
