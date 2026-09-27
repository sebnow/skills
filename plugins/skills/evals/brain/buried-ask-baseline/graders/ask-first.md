---
type: llm
focus: last_message
---
The user is working with a decision-making agent that has just received the
last of several worker reports full of file references, requirement ids and
regulation citations. The reply exists to get decisions from the user. A
reader should learn from the opening alone what the situation is and what
they are being asked to do; a reply that opens with findings, file maps or a
"here's what's there" section makes the reader dig for the ask.

PASS: the first two sentences of the reply (ignoring a heading line) together
state what is going on in plain words and what the user is being asked to
decide or answer, or say plainly that nothing is needed from the user yet.

FAIL: the first two sentences do not say what the user is being asked, or
the ask first appears only after a block of context, findings or evidence
(for example after a "What's there" or "Findings" section, or after bullet
lists of code references). A generic opener such as "I have what I need,
here's the frame" does not count as stating the ask.
