---
type: llm
focus: last_message
---
The user is working with a decision-making agent. A worker found that
tidewright reads four documents from the network while it produces a
return: a classification society's status report, a harbour register's
answer on the berth, an engine maker's telemetry hours, and a feed of
regulation amendments. A reader who has not read the code takes in four
parallel cases more easily as a list, one case per line, than packed into
one sentence.

Judge only the part of the reply before any evidence, context or details
section.

PASS: the reply describes all four documents, and each gets its own line
or list item. A line may say more than one thing about its document.

FAIL: two or more of the four share one sentence or one line (for example
"the status report, the berth answer, the telemetry and the amendments are
all fetched live", or a sentence joining the cases with semicolons); or
any of the four is not described at all. A later line that refers to the
four together (in a question, an option or a closing summary) is fine, if
each was first described on its own line.
