---
type: llm
focus: last_message
weight: 1
---
The input is a research report compiled by a worker agent: findings quoted from web sources with URLs, a sources list, notes on how it was gathered (who compiled it, that pages were fetched with WebFetch, that quotes marked [s] came from the fetch tool's summariser), the worker's own UNVERIFIED inferences, and open questions ("The skill must choose."). A skill directory is loaded by an agent doing the skill's task, here writing an ADR. The report, and any account of how the rules were researched, does not help that agent write an ADR; where a rule's source matters, it belongs in the commit message that adds the rule.

Judge only the files the reply proposes for skills/adr-writing/. Commentary in the reply outside those files does not count, nor do files proposed elsewhere (for example eval cases).

PASS if every proposed file in the skill directory holds guidance an agent applies while writing or updating an ADR, and none of them carries the report or its process narrative. A reference file of task guidance (for example a template, or a checklist for deciding whether a decision needs an ADR) is fine.

FAIL if any proposed file in the skill directory copies or digests the report as research, for example: a research, evidence, sources, or bibliography file; findings attributed to named sources or with URLs as the justification for a rule; a note on who compiled the material or with which tool; [s] or UNVERIFIED markers; the worker's inferences; or open questions. Also FAIL if the reply proposes no files for the skill directory, since there is nothing to judge.
