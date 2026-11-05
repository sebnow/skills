---
type: regex
target: last_message
pattern: '^(?:(?!\n(?:#{1,6}[ \t]+|\*\*)[^\n]{0,30}(?:evidence|context|details)[^\n]{0,30}\n)[\s\S])*?(?:(?=[^\s]*\w)(?![^\s]*[.!?][)"''*_”’\]]*(?=\s|$))[^\s]+[ \t]+(?:[^\w\s]+[ \t]+)*){25}[^\s]*\w'
flags: i
match: not_contains
---
Fails when a sentence above the evidence heading (a markdown heading or
bold line naming evidence, context or details) runs past 25 words, the
limit brain.md takes from the GOV.UK style guide. A word is a
space-separated token holding a letter or digit, so list markers and
dashes do not count. A sentence ends at ".", "!" or "?" before a space or
line end (closing quotes, brackets and bold markers allowed after it) or
at a line break; a colon or semicolon does not end it. "e.g." ends a
sentence here, so the check errs towards passing. Without an evidence
heading it reads the whole reply.
