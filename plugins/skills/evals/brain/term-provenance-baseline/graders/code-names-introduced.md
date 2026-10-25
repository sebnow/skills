---
type: llm
focus: trace
---
The trace may run past the agent's reply to the scout's report: a worker
it dispatched can report back, and the agent then replies again. Judge only
the agent's first reply to the report, up to the point where it first ended
its turn; ignore everything after that, the steps a dispatched worker took
itself, and any brief the agent sent a worker.
A <task-notification> marks a worker reporting later: everything from the
first one on is outside the first reply, however many questions it asks.

The user owns the repository but has not read the scout's report, which
names code and requirement identifiers: `voyage.Dedupe` (the function that
drops a voyage overlapping an earlier one from another import), `Clash`
(the record of two overlapping voyages), `Voyage.Import` (which import a
voyage came from), `kept`, `From`/`To`, and REQ-0108 (the requirement that
overlapping voyages are an open finding and neither is dropped).

PASS: every one of these identifiers the reply uses above its evidence or
context section is said, where it first appears, to be what it is in plain
words (for example "the dedupe step, `voyage.Dedupe`, which drops the later
voyage"). Identifiers that appear only in an evidence or context section
at the end need nothing. A reply that uses none of them passes.

FAIL: the reply uses one of these identifiers above its evidence section
with no word on what it is or does (for example "REQ-0108 contradicts
Dedupe", "change Clash to keep both").
