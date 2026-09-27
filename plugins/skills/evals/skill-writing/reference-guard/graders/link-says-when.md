---
type: llm
focus: last_message
weight: 1
---
The reply splits a skill: a shorter SKILL.md and a new reference file holding the rate limit details. The reference file is only loaded when the agent decides to read it, and it decides from SKILL.md alone, so the condition for reading the reference must sit in SKILL.md at the link.

PASS if the SKILL.md in the reply links the reference file and states, at that link, the concrete situation in which to read it (for example: when handling a 429, when setting retry or backoff for a partner call, when writing a batch call, or when the partner is Northwind or Acme).

FAIL if the link in SKILL.md gives no situation ("see references/rate-limits.md for details", "rate limit specifics are in references/rate-limits.md"), or if the only statement of when to read the file is inside the reference file itself, or if SKILL.md has no link to the reference file.
