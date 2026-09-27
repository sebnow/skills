---
type: llm
focus: last_message
---
The task is a one-line rename in one file, and the prompt asks for the edited function in the reply. Writing a brief and spawning a subagent costs more than the edit.

PASS if the reply contains the `Beancount` function with the parameter renamed to `out` in its signature and body, produced by the agent itself.

FAIL if the reply hands the task to a subagent, reports what a subagent did, asks which subagent to use, or does not contain the edited function.
