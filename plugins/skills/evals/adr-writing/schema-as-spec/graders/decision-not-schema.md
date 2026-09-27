---
type: llm
focus: last_message
weight: 1
---
The workspace holds a migration with the full DDL of an `events` table (16 columns, 3 check constraints, 6 named indexes) and a design note saying the team chose one table with a JSONB `payload` column over one table per event subtype. The note gives the reasons (27 types growing by about two a month, a migration per type, a 27-way UNION for the feed; filtering only on real columns) and a boundary: anything needing a foreign key into the payload, or queried by a payload field on a request path, gets promoted to a real column or its own table. The migration is the source of truth for the schema; the ADR should record the decision, not restate the schema.

Judge only the ADR in the reply.

PASS if the ADR explains why JSONB won over per-subtype tables and states the promotion boundary (foreign key or request-path query on a payload field means promote), and any schema it shows is a small illustrative shape (a handful of columns at most, or a sentence naming which fields are real columns) rather than the table definition.

FAIL if the ADR reproduces the DDL or most of it, lists the table's columns one by one (more than about six), lists the indexes or constraints by name or definition, or omits either the reason for choosing JSONB or the promotion boundary.
