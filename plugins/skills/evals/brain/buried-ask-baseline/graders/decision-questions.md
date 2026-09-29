---
type: llm
focus: last_message
---
The reply asks the user for decisions in a numbered list. A reader must be
able to tell from each item what is being decided. An item whose headline is
the agent's recommendation stated as a fact ("How rates enter the record:
fetch from the API, stored raw") leaves the reader unable to tell what the
decision is or what the alternatives are, even when the reply later says
"Recommendation:".

PASS: every numbered decision item (the list before the context or evidence
section; ignore a numbered restatement at the close) puts to the user a
question they can answer, and that question says what is being decided. The
item's bold headline or first sentence is the question or leads into it, not
the recommendation.

FAIL: any decision item has no question to the user, or its headline or
first sentence states the recommended outcome as a settled fact, or its only
question is a generic "agree?" or "confirm?" that does not say what is being
decided.
