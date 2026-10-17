---
type: llm
focus: last_message
weight: 1
---
The report is an architecture review of a Go module. It should present each architectural transformation it proposes as a choice between designs, so the reader can compare them and pick one.

PASS if every transformation the report proposes offers at least two alternative designs, and states pros and cons (or trade-offs) for the alternatives. Doing nothing counts as an alternative only if its trade-offs are stated. Defects or observations listed apart from the proposed transformations are not transformations and are not judged here.

FAIL if any proposed transformation is given as a single design without an alternative, or if alternatives are listed without their pros and cons, or if the report proposes no transformation at all.
