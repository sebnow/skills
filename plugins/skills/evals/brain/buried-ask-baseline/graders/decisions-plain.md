---
type: regex
target: last_message
pattern: '^(?:(?!REQ-\d|ADR-\d|[\w/-]+\.[a-z]{1,4}:\d|\breg(?:ulation)?\.?[ \t]*\d|§)[\s\S])*?\n[ \t]*(?:\*\*)?1[.)][ \t](?:(?!REQ-\d|ADR-\d|[\w/-]+\.[a-z]{1,4}:\d|\breg(?:ulation)?\.?[ \t]*\d|§)[\s\S])*\n(?:#{1,6}[ \t]+|\*\*)[^\n]{0,30}(?:evidence|context|details)[^\n]{0,30}\n'
flags: i
match: contains
---
Passes when the reply has a numbered list (a line starting "1." or "**1.")
followed by a heading line naming evidence, context or details, and nothing
before that heading carries a `path:line` reference, a `REQ-` or `ADR-` id,
or a regulation citation ("reg. 11(5)", "regulation 17", "§4"). So the lead
and the numbered decisions are in plain words and every citation sits in the
evidence section. Replaces an llm rubric whose one-word judges failed a reply
with no citation in any decision item (round one, with-section run 2), and
flipped to PASS when asked to name the offending item.
