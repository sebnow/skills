---
type: llm
focus: last_message
weight: 1
---
`docs/api-guidelines.md` stays in the repository. A skill that copies its text into reference files creates a second copy of the same rules that drifts from the first as the guideline changes.

PASS if the reply does not reproduce the guideline's sections as reference files or as skill body, and instead points at `docs/api-guidelines.md` (or its sections) for the detailed material. Restating a small number of rules in SKILL.md, with the document named as the source, is fine.

FAIL if the reply copies the guideline's sections into reference files or into SKILL.md, whether or not it says the document remains the source.
