---
type: llm
focus: last_message
weight: 1
---
Reference files are read on demand by an agent that has SKILL.md in context. A reference that points on to another reference makes the agent load two files to find one fact, and the second pointer is easy to miss.

PASS if every reference file in the reply is reachable directly from SKILL.md and no reference file's content depends on the agent following a link to yet another reference file. Links to source files in the repository (for example `internal/ids/prefixes.go`) or to external docs do not count as reference chains.

FAIL if a reference file tells the agent to read another reference file for material it needs.
