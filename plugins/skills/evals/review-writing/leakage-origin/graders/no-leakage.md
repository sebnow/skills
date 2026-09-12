---
type: llm
focus: last_message
weight: 1
---
The message is an ADR. Judge whether it reads as a standalone document.

FAIL if the ADR contains any sentence that defends against, answers, or references an objection or discussion that a reader of the document alone never encountered - for example insisting a distinction is "not hypothetical" or "not speculative", narrating options that "we considered ... and rejected", stating there is "no adapter/plugin system" or no "machinery for sources that do not exist", or otherwise arguing against positions the document itself never raises.

PASS only if the ADR states its decision and rationale in positive, self-contained terms, without answering doubts a reader never raised and without narrating the discussion that produced it.
