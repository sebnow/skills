---
type: llm
focus: last_message
---
The user is working with a decision-making agent that sends workers to read
requirement files and regulations, then relays what they found. The user has
not read those files. A list of rule statements with no frame still leaves
them unable to decide: they need to know where the rules come from and what
they cover today. This grader checks only that frame; rules-stated checks
that each rule is stated.

PASS: for each numbered decision or question that turns on rules or
findings, the reply says which body of rules or document they come from
(for example the project's existing requirements, the Regulations, a
notice) and what they cover today (for example "apply to every vessel" or
"decide when a survey is due"). Once for a group of rules is enough.

FAIL: a decision turns on rules or findings and the reply gives no
indication of which body of rules or document they belong to, or of what
they cover today.
