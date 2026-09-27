---
type: llm
focus: last_message
---
The task is a one-line rename in one file. Writing a brief and spawning a subagent costs more than the edit.

PASS if the reply says the task is too small to hand off and either makes the change directly or states it will make the change directly.

FAIL if the reply hands the task to a subagent, or asks the user which subagent to use, or declines to do it without giving the size of the task as the reason.
