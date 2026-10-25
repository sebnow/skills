---
type: llm
focus: trace
---
The trace may run past the agent's reply to the worker's report: a worker
it dispatched can report back, and the agent then replies again. Judge only
the agent's first reply to the report, up to the point where it first ended
its turn; ignore everything after that, and ignore the steps a dispatched
worker took itself.
A <task-notification> marks a worker reporting later: everything from the
first one on is outside the first reply, however many questions it asks. A
brief the agent sends a worker is not a question to the user.

The session was compacted. The summary the agent works from says the
conventions are "one decision per ADR; ADRs supersede whole, never amended;
requirements withdrawn and replaced, never edited", and points at
docs/conventions.md, which the agent read before compaction and can read
again. That file says: (rule 2) an accepted ADR is never edited or amended in part; a change
is a new ADR that supersedes it whole; (rule 3) when a decision changes,
every accepted ADR that states the old decision is superseded in the same
change; (rule 6) a requirement whose rule changes is never edited: it is
marked withdrawn with a "Replaced by" line and the new rule goes in a new
file. A worker then reported two open points: ADR-0004 still says elections
are per operator (superseding ADR or leave it?), and REQ-0139 and REQ-0140
still describe the operator's election (edit them or leave them?). The
conventions answer both: ADR-0004 gets a superseding ADR; REQ-0139 and
REQ-0140 are withdrawn and replaced by the new requirements. The user told
the agent to follow the conventions, and when the agent earlier asked them
to confirm a point the conventions settle, answered "yes, obviously, that's
what the conventions say" (quoted in the summary).

PASS: the reply asks the user nothing about either point. It may say what
will be done about them (supersede ADR-0004, withdraw REQ-0139 and
REQ-0140), with or without citing the conventions, or not mention them.

FAIL: the reply puts either point to the user as a question or decision in
any form: a numbered question, "do you confirm", "OK to supersede ADR-0004?",
"shall I withdraw them?", options to pick from, or a recommendation the
user is asked to approve.
