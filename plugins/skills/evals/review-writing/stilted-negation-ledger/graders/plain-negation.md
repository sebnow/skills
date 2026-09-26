---
type: llm
focus: last_message
weight: 1
---
The message contains a revised document. Judge only its phrasing of absence and negation.

Judge only the final revised document (the last fenced markdown block). Ignore commentary outside it, including quotes of the original draft.

FAIL if the revised document still states an absence in a compressed, aphoristic form where an ordinary writer would use a plain verb negation. Examples of the pattern: "X has no Y", "A and B share no Z", "X carries no Y", "X holds no Y", "leave no trace of", and riddle-like setups such as "has to generate what the model does not store:" that withhold the subject until after a colon. Plain forms such as "X does not have Y", "A and B do not share Z", or "The model does not store X, so the export generates it" are acceptable.

Check each quoted example against the revised document; if any of them, or a sentence built the same way, remains, FAIL.

PASS only if every absence in the revised document is stated plainly.
