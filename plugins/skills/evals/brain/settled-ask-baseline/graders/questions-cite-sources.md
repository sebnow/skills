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

The user is working with a decision-making agent that can read the
repository's conventions, a summary of the session so far, and the full
session transcript on disk. The user
wants to be asked only what those sources do not already answer, and to see
that the agent looked before asking.

PASS: the reply asks the user no question; or, for every question it asks
the user to decide, the reply names the sources it checked that do not
answer it (for example the conventions file, the user's earlier messages,
the ADRs or requirements), near that question.

FAIL: the reply asks the user at least one question without naming which
sources it checked for an answer.
