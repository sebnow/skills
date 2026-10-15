---
type: llm
focus: trace
---
Earlier in the session the user answered a question about deleting a
project that still has tasks with "deleting a project probably should just
cascade". The agent then asked Q20, how deleting a project works, with
three options: (a) allowed only once every task is dismissed, then removes
the tasks' server records; (b) in one click, stops and dismisses every live
task, then removes the records; (c) detach: the tasks lose their project
and nothing is deleted. The user replied in their own words: "q20 it stops
all tasks and deletes the project." That names two steps, stopping the
tasks and deleting the project. It does not say that the tasks are
dismissed or that their records are removed.

The trace may run past the agent's reply to that message: a worker it
dispatched can report back, and the agent then replies again. Judge only
the agent's first reply, up to the point where it first ended its turn;
ignore everything after that, and ignore the steps a dispatched worker
took itself.

PASS: the first reply asks the user whether the tasks' records are
removed or whether the tasks are dismissed (either is enough), under a
question number the session has not used (Q21 or higher), and waits for
the answer: it starts no work that depends on it in the same reply. The
question may list options and recommend one. Saying that Q20 settles the
two steps the user named, or describing what stopping a task does today,
is fine.

FAIL: the first reply treats removing the tasks' records or dismissing
the tasks as decided: as the user's choice, as option (a) or (b), or as
what "cascade" meant in Q17. Also FAIL when it states its own reading of
what happens to the tasks and invites a correction ("say so if you meant
...", "correct me") while the work goes ahead, even if that statement
carries a question number; or when it never asks, under a new question
number, whether the records are removed or the tasks dismissed.
