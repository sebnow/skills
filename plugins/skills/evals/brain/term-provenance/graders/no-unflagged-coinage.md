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

In past sessions the user complained that the agent invents names for things ("Entry is
too broad of a term", "it's not clear what that refers to") and then uses
them as if they were established. A coined name is a label the reply gives
to a concept, set off as a term (in quotes, bold or italics, capitalised as
a name, or used repeatedly as a noun for a specific design element), that
does not appear in the user's messages, the scout's report, or the
project's own vocabulary (voyage, import, clash, finding, open finding,
engine hours, return, ledger, evidence, election, vessel, operator).
The report's own terms (track fusion, leg fingerprint, primary device, gap
fill, merge candidate, twin badge, HDOP gating, the 90-second rule) are not
coinages.

PASS: the reply introduces no coined name; or each coined name is marked as
the agent's own proposal where it first appears (for example "call this,
for now, a 'merged voyage'; tell me if you want another name") and its
meaning is stated.

FAIL: the reply uses a coined name as if it were established, without
saying it is the agent's own and what it means (for example "each voyage
gets a source set", "the canonical track wins", "a twin voyage" as
standing terms).
