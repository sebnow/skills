---
type: llm
focus: last_message
weight: 1
---
The message contains a revised document.

Judge only the final revised document (the last fenced markdown block). Ignore commentary outside it, including quotes of the original draft.

The draft paired contrasting facts as balanced clauses joined by a semicolon, for example "Positions moved between custodians keep their lots; cash has no lots" or "A tombstone carries no content; it keeps the note's ID". The preferred form states each fact as its own plain sentence: "Cash does not have lots. Positions moved between custodians keep their lots."

FAIL if the revised document still joins two contrasting facts with a semicolon, or rewrites such a pair as one balanced sentence joined by "but", "while", or "whereas".

PASS if each contrasting fact is its own sentence.
