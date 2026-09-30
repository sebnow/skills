---
type: regex
target: last_message
pattern: '^(?:(?!\n(?:#{1,6}[ \t]+)?(?:\*\*)?\d{1,2}[.)](?=[ \t*]))[\s\S])*?(?:\n(?:#{1,6}[ \t]+)?(?:\*\*)?\d{1,2}[.)](?=[ \t*])(?:(?!\n(?:#{1,6}[ \t]+)?(?:\*\*)?\d{1,2}[.)](?=[ \t*])|\n(?:[ \t]*\n)?[ \t]*(?:[-*>][ \t]+)?(?:\*\*|__|_|\*)?(?:I |my )?recommend|\n(?:#{1,6}[ \t]+|\*\*)[^\n]{0,30}(?:evidence|context|details)[^\n]{0,30}\n)[\s\S])*?\?[^?\n]*(?:(?!\n(?:#{1,6}[ \t]+)?(?:\*\*)?\d{1,2}[.)](?=[ \t*]))\n[ \t]*[^\s?][^?\n]*)*?(?:\n[ \t]*\n(?:[ \t]*[^\s?][^?\n]*:[ \t]*\n(?:[ \t]*\n)?)?[ \t]*(?:[a-z][.)]|\([a-z]\)|[-*+])[ \t][^?\n]*(?:\n(?:[ \t]*\n)?[ \t]*(?:[a-z][.)]|\([a-z]\)|[-*+])[ \t][^?\n]*|\n[ \t]*[^\s?][^?\n]*)*?)?\n(?:[ \t]*\n)?[ \t]*(?:[-*>][ \t]+)?(?:\*\*|__|_|\*)?(?:I |my )?recommend(?:(?!\n(?:#{1,6}[ \t]+)?(?:\*\*)?\d{1,2}[.)](?=[ \t*])|\n(?:#{1,6}[ \t]+|\*\*)[^\n]{0,30}(?:evidence|context|details)[^\n]{0,30}\n)[\s\S])*)+\n(?:#{1,6}[ \t]+|\*\*)[^\n]{0,30}(?:evidence|context|details)[^\n]{0,30}\n'
flags: i
match: contains
---
Passes when every numbered decision item before the evidence heading (a
markdown heading or bold line naming evidence, context or details) has a
question, and a line starting "Recommend" (optionally bold, bulleted or
"I recommend") follows it: after the question's last "?" only lines without
a "?" and at most one blank line come before the recommendation, except
for one list of options ("a.", "(a)" or a bullet, optionally under a line
ending in ":") set off by a blank line. Fails an item with no "?", a
recommendation on the same line as the question ("Q? Recommend fetch."), a
recommendation split from its question by another paragraph, or an item
with no recommendation line. Items are lines starting
"1.", "**1." or "### 1."; nested lists must be indented. Does not judge
whether the question is a good one; decision-questions does that.
