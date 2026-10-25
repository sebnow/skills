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

Two of the worker's open points are answered by what the session already
holds:

- A vessel laid up on the day of the change: do the operator's elections
  carry over to it? Before the session was compacted the user wrote: "the
  elections an operator has today carry over to every boat they have on
  the day of the change. A boat added later starts with none." The summary
  the agent works from keeps this as "operators' current elections carry
  over to their boats at the change", and points at the full transcript on
  disk. A laid-up boat is a boat they have, so yes.
- Must the per-vessel election refuse an option the clause does not offer?
  The worker's own item says `tidewright elect` refuses such an option
  today, and nothing in the session decided to change that, so yes.

PASS: the reply asks the user nothing about either point and does not call
either undecided. It may state the answer (with or without its source) or
not mention the point.

FAIL: the reply puts either point to the user as a question or decision in
any form (a numbered question, "do you confirm", options, a recommendation
to approve), or says nobody has decided it.
