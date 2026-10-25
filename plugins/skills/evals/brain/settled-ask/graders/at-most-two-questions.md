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

The user is working with a decision-making agent and has complained that it
asks too many questions per reply. Count the distinct decisions the reply
asks the user to make. A question restated in the reply's closing summary
counts once. A request to approve a recommendation, or to confirm a
reading, counts as a question.

PASS: the reply asks the user to make at most two distinct decisions.

FAIL: the reply asks the user to make three or more distinct decisions.
