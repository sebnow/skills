---
type: llm
focus: last_message
---
The user is working with a decision-making agent that sends workers to read
requirement files, regulations and source documents, then relays what they
found. The user has read none of those files or the worker reports; they see
only this reply. A decision about rules the user has never seen is one they
cannot make, however well the reply is laid out (an ask-first reply with
numbered decisions and an evidence section can still fail here), so each
rule or finding a decision turns on must be stated in the reply itself.

PASS: for every numbered decision or question the reply asks the user to
answer, each rule, requirement or finding it concerns is stated in the reply
in plain words, saying what the rule requires or what the finding is, so a
reader who never opened the source files could answer. A group of rules
passes only when each rule in the group is stated (for example "an engine
replaced mid-season starts again from zero hours; one past its overhaul
limit is flagged" rather than "the overhaul rules"). Stating them in a
section before the decision is fine.

FAIL: any decision or question refers to rules, requirements or findings by
count, title, identifier or section name only ("the five overhaul rules",
"the three survey rules", "REQ-0124", "§4 and §9", "the Schedule 2 table"),
or its answer depends on what a file, notice or section says without the
reply saying it.
