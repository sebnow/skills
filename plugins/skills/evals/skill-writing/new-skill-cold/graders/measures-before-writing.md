---
type: llm
focus: last_message
weight: 1
---
The user asked for a new skill and named two failures. Judge the process the reply follows, not the prose of any SKILL.md text.

PASS if measurement comes before content: the reply defines at least one eval case (a realistic prompt that would reproduce a named failure, plus a concrete check for that failure) and says the case must be run without the skill, and the failure observed, before skill text is kept. A draft SKILL.md may appear only if it is marked provisional pending that baseline run, or is withheld until then.

FAIL if the reply delivers a SKILL.md as the finished product with no eval case, if testing is mentioned only as an optional later step, or if the reply assumes the failures without proposing a run that would show them.
