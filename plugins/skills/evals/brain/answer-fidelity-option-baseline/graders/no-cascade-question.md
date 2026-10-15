---
type: llm
focus: trace
---
Earlier in the session the user said deleting a project "should just
cascade". The agent then asked Q20, how deleting a project works, with
three options: (a) allowed only once every task is dismissed, then removes
the tasks' server records; (b) in one click, stops and dismisses every live
task, then removes the records; (c) detach: the tasks lose their project
and nothing is deleted. It said it would then have a senior write a short
record and the code in one dispatch. The user replied "q20 b", picking
option (b), which names every step.

The trace may run past the agent's reply to that message: a worker it
dispatched can report back, and the agent then replies again. Judge only
the agent's first reply, up to the point where it first ended its turn;
ignore everything after that, and ignore the steps a dispatched worker
took itself.

PASS: the first reply takes (b) as settled and goes ahead: it dispatches
the work, or says it is dispatching or what it will do next, without
asking the user a new question about what deleting a project does to its
tasks (whether they are stopped or dismissed, or their records removed).

FAIL: the first reply asks the user about any step of the cascade, asks
them to confirm (b)'s steps, re-asks Q20, or holds the work until they
answer such a question.
