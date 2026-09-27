---
type: llm
focus: last_message
weight: 1
---
The reply proposes a SKILL.md for skills/adr-writing/, written from a research report. SKILL.md is what the agent reads first; a pointer from it to the research sends the agent to material about where the skill's rules came from rather than how to write an ADR.

Judge only the text inside the proposed SKILL.md. Ignore the reply's commentary around the files, including remarks about the research report, where it stays, or what was left out of the skill.

PASS if no sentence in SKILL.md directs the reader to a file that holds research, evidence, or sources for the skill's own rules, or to the report itself (research/adr-writing-2026-09.md), whether by markdown link or by naming the file in prose. These are task guidance and do not count as such a pointer: telling the ADR author where to get a decision's rationale ("sources are the user, issues, design docs"); a link to a template or checklist file; a template section for links to design docs or related ADRs.

FAIL if SKILL.md links or names such a file, for example "Read references/research.md when a rule here seems wrong or you want its source", "Sources: references/sources.md", or "The research behind these rules is in research/adr-writing-2026-09.md". Also FAIL if the reply has no SKILL.md.
