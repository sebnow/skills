---
type: llm
focus: last_message
weight: 1
---
The input is a guideline document of about 500 lines, `docs/api-guidelines.md`, covering many sub-areas of HTTP API design, and it stays in the repository. Skill text is loaded into the agent's context in full whenever the skill triggers, while anything the skill points at is loaded only when the agent decides to read it. A skill that reproduces the whole document loads hundreds of lines for every endpoint change, most of them irrelevant to the change at hand.

PASS if the SKILL.md in the reply is clearly shorter than the document, holds only what every endpoint change needs, and makes the rest reachable through pointers that each name the situation in which to follow them. The pointer may go to the document itself, to named sections of it, or to reference files; what matters is that the agent can tell from SKILL.md when to read further.

FAIL if SKILL.md reproduces substantially all of the guideline, or if the pointers to further material give no situation ("see docs/api-guidelines.md for the rest", "see references/errors.md for details"), or if sub-areas of the guideline are neither in SKILL.md nor reachable from it.
