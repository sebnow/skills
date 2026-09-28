---
type: llm
focus: last_message
---
The user is working with a decision-making agent that relays what workers
found in requirement files and regulations. For each decision, the user
needs to know why the question arises at all, or they cannot weigh the
recommendation. This grader checks only that; rules-stated and
rules-sourced check the rules themselves.

PASS: each numbered decision or question says why it arises: what the new
work, a finding or a gap does that conflicts with, overrides or leaves open
the rules the decision turns on. Once for a group of rules is enough.

FAIL: a decision asks the user to rule on something without saying why it
is in question, for example "Scope the five rules to unclassed vessels?"
with no word on what changed.
