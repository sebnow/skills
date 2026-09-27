---
type: llm
focus: last_message
---
The reply should carry a decision layer the user can act on without reading
code: numbered decisions, each with the agent's recommendation in plain
words. Evidence such as file paths with line numbers (`flow.go:74`),
requirement ids (`REQ-0114`) and regulation citations (`reg. 22(6)`,
`CCVR 2021 reg. 11(5)`) belongs in a separate evidence section, not inside
the decision items. Identifier names in backticks (a type or command name)
are allowed; paths with line numbers, REQ ids and regulation clause
citations are not.

PASS: the reply contains a numbered list of decisions (or questions for the
user), each item gives a recommendation (what the agent would do or
suggests), and no item, including its sub-bullets, contains a `path:line`
reference, a `REQ-` id, or a regulation citation such as "reg. 4(1)".

FAIL: there is no numbered list of decisions; or any decision item lacks a
recommendation; or any decision item or its sub-bullets contains a
`path:line` reference, a `REQ-` id or a regulation citation. Lettered lists
(A/B/C) count as numbered.
