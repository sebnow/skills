---
type: llm
focus: trace
weight: 1
---
The repository's `docs/adr/0003-service-naming.md` is an ADR with status Accepted that defines a naming scheme for workloads. The user asked to add a new rule (a `svc-` prefix for Kubernetes Service objects) to that ADR. An accepted ADR records a decision as it was made; a new rule is a new decision and belongs in a new ADR (which may amend or supersede ADR-0003) or in a living conventions document, not in an edit to the accepted record. The agent had no edit tools, so judge what it proposes, not what it could do. Judge every message the agent wrote in the trace, not only the last one. Reading ADR-0003 with a tool is not an edit.

PASS if the agent declines to put the new rule into ADR-0003, explains that ADR-0003 is accepted and should not be rewritten, proposes a new ADR (for example ADR-0004 amending ADR-0003) or a living document for the rule, and nowhere offers or proposes any other change to ADR-0003. The one exception is a status-line pointer on ADR-0003 to the new ADR ("Superseded by ADR-0004" or "Amended by ADR-0004"): suggesting that still passes, because it adds a pointer and not the rule.

FAIL if any message edits ADR-0003 to include the rule, presents a revised ADR-0003 containing the rule, gives a diff or replacement text adding the rule to ADR-0003, or offers to edit ADR-0003 in any way beyond that status-line pointer ("want me to update 0003 too?", "I can also add a note to 0003's Decision section", a "Related" or "See also" link in its body). Such an offer fails the reply even when it comes alongside a correct proposal for a new ADR.
