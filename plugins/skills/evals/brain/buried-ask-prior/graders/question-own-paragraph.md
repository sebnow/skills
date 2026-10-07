---
type: regex
target: last_message
pattern: '\n(?:#{1,6}[ \t]+)?(?:\*\*)?\d{1,2}[.)][ \t](?:(?!\n(?:#{1,6}[ \t]+|\*\*)[^\n]{0,30}(?:evidence|context|details)[^\n]{0,30}\n)[\s\S])*?[a-z)"''”’][.:](?:\*\*)?[ \t]+[^\n.?]*\?'
flags: i
match: not_contains
---
Fails when, in the decisions before the evidence heading, a question ends a
line that also holds an earlier sentence, as in "...passing output. Do we
adopt the ladder?". The question should start its own paragraph so it reads
apart from the context. A question containing a "." (a file name) is not
caught, and the item number's own "." does not count as a sentence end.
