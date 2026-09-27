---
type: llm
focus: last_message
weight: 1
---
The repository's `docs/adr/` holds three ADRs, `0001-record-architecture-decisions.md`, `0002-postgresql-for-ledger-storage.md` and `0003-idempotency-keys-on-write-endpoints.md`. Each has exactly this layout: a `# Title` line with no number or prefix, then `## Status` (value "Accepted"), `## Context`, `## Decision` and `## Consequences`, in that order, with prose paragraphs and no other sections. ADR 0001 states this layout as the convention. The user asked for the next ADR.

Judge only the new ADR in the reply.

PASS if the new ADR has a `# Title` line in the same style (no "ADR-0004:" or "4." prefix) followed by exactly the four sections `## Status`, `## Context`, `## Decision`, `## Consequences` in that order, and the proposed filename follows `NNNN-short-title.md` with number 0004.

FAIL if it adds a top-level section the existing ADRs do not have (such as "Alternatives Considered", "Decision Drivers", "Options" or "References"), drops, renames or reorders any of the four, adds a metadata block or table the others lack, prefixes the title with a number, or proposes a filename outside the `0004-short-title.md` pattern.
