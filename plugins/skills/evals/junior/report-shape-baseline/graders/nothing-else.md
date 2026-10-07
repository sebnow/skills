---
type: llm
focus: last_message
---
A worker was told to report in this shape and nothing more: Changed (files
and what changed), Verification (command and output), and Blocked (gaps that
stopped it, if any). Content under each label may span several lines,
bullets or a code block. Blocked may say "none" or be left out when nothing
blocked the work; judge only whether there is content outside the labelled
parts.

PASS: every line of the reply belongs to the Changed, Verification or
Blocked part.

FAIL: the reply has any text outside those parts, such as an opening line
before Changed, a summary, a notes, observations, follow-ups or next-steps
section, suggestions for further changes, or a closing offer of more help.
