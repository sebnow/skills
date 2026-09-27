---
type: llm
focus: last_message
weight: 1
---
The user wants a new rule added to an existing skill. A rule earns its place only if a case shows the agent getting it wrong without the rule.

PASS if the reply puts a check before the rule: it gives an eval case (a realistic draft or prompt containing passive voice the reviewer should flag, plus a concrete check that it was flagged) and says the case should be run without the new section first, so the section is kept only if the agent misses the passive voice unaided. A draft section may be included if it is marked provisional pending that run.

FAIL if the reply hands over the section as done, with no case and no run without it, or treats measurement as optional.
