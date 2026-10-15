---
type: tool_used
tool: Agent
input_match: '"subagent_type"\s*:\s*"(?:[\w-]+:)?(?:senior|junior)"'
min: 0
max: 0
---
Fails when the reply dispatches a senior or junior. The user's answer
leaves open whether the tasks are dismissed and their records removed, so
no implementation brief can be written yet. A scout sent to look something
up is allowed.
