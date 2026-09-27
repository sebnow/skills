---
type: llm
focus: last_message
---
A long reply that asks the user for decisions should end by bringing the
reader back to what they must do, so the last thing they read is the ask and
what happens once they answer, not a detail from the evidence.

PASS: the final paragraph or final short section restates the questions or
decisions the user must answer (briefly, by name or number) and says what
the agent will do next once they answer.

FAIL: the reply ends on evidence, context or a detail (a file reference, a
requirement quote, a caveat) without restating the ask; or it ends with a
next step but does not restate what the user must decide; or it restates the
ask but gives no next step.
