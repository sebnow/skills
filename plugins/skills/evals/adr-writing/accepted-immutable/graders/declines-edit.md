---
type: llm
focus: last_message
weight: 1
---
The repository's `docs/adr/0003-service-naming.md` is an ADR with status Accepted that defines a naming scheme for workloads. The user asked to add a new rule (a `svc-` prefix for Kubernetes Service objects) to that ADR. An accepted ADR records a decision as it was made; a new rule is a new decision and belongs in a new ADR (which may amend or supersede ADR-0003) or in a living conventions document, not in an edit to the accepted record. The agent had no edit tools, so judge what it proposes, not what it could do.

PASS if the reply declines to put the new rule into ADR-0003, explains that ADR-0003 is accepted and should not be rewritten, and proposes a new ADR (for example ADR-0004 amending ADR-0003) or a living document for the rule. Also suggesting a one-line status note on ADR-0003 that points to the new ADR ("Amended by ADR-0004") still passes, because it adds a pointer and not the rule.

FAIL if the reply edits ADR-0003 to include the rule, presents a revised ADR-0003 containing the rule, gives a diff or replacement text adding the rule to ADR-0003, or offers to do so ("want me to update 0003?"), even alongside other suggestions.
