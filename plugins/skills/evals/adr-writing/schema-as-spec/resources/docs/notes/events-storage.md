# Events storage - design session notes (2026-09-10)

Attendees: Marta, Idris, Chen

- One `events` table with a JSONB `payload`, not a table per subtype
  (`invoice_events`, `login_events`, ...). 27 event types today, product is
  adding roughly two a month. A table per subtype means a migration and a new
  repository type for every new event, and the activity feed would need a
  27-way UNION.
- The feed and the audit export only filter on tenant, time, type, subject and
  correlation id. Those are real columns. Nothing reads payload fields in a
  WHERE clause today except the reconciliation job, which reads
  `payment.captured` amounts once a night.
- Payload shape per type lives in `internal/events/schemas/*.json` and is
  validated on write in the Go service, not by the database. `payload_schema_rev`
  lets us change a shape without rewriting old rows.
- Line we agreed: if something needs a foreign key into the payload, or gets
  queried by a payload field on a request path, it gets promoted to a real
  column or its own table. Reconciliation is the first candidate if it moves
  off the nightly batch.
- Chen worried about losing type safety in SQL. Accepted: validation happens in
  one place in the service, and the GIN index covers ad-hoc support queries.
