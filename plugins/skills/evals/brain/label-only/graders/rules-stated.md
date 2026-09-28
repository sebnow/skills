---
type: llm
focus: last_message
---
The user is working with a decision-making agent that sends workers to read
requirement files and regulations, then relays what they found. The user has
not read those files or the worker reports; they see only this reply. A
decision about rules they have never seen is one they cannot make, however
well the reply is laid out, so the reply must say what each rule says. This
grader checks only that; rules-sourced and rules-bearing check the frame.

PASS: every rule or finding that a numbered decision or question turns on
has its content stated in the reply in plain words: what the rule requires
or what the finding is. A decision to change, scope, reword, point at or
keep a group of rules turns on those rules, even when the edit itself is
small. A clause per rule is enough, and a group passes
when each of its rules is stated this way, even in one sentence (for
example "an engine replaced mid-season starts again from zero; one past its
overhaul limit is flagged").

FAIL: a decision refers to rules or findings by count, title or identifier
("the five overhaul rules", "REQ-0124", "§4 and §9"), or by subject without
what they require ("the hours-based overhaul rules", "the passenger,
hull-work and laid-up rules"), so the reader cannot tell what the rules
say. Saying that rules exist, stop applying or need a pointer is not saying
what they require.
