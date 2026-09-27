---
type: llm
focus: last_message
weight: 1
---
The conversation before this ADR chose a pg_trgm GIN index over OpenSearch for catalog-api's suggest endpoint. Two facts decided it. A k6 load test (script at `loadtest/suggest.js`, 200 rps for 10 minutes against staging) measured p99 340ms for the current ILIKE query and 90ms with the trigram index, against a 150ms p99 SLO. And the platform team ruled out adding a second datastore this quarter because of an RDS migration; the user learned this in a Slack thread. The ticket is CAT-412. A reader of the ADR a year from now has no access to that Slack thread or to the chat that produced the ADR.

Judge only the ADR in the reply.

PASS if the ADR body states the measured latency result (at least the 90ms p99 with the trigram index, set against either the 340ms baseline or the 150ms SLO) and states the constraint that no second datastore can be added this quarter, as facts in the ADR itself; and every source it cites is an in-repo path (such as `loadtest/suggest.js`), a commit id, or a ticket id (such as CAT-412). Attributing the constraint to the platform team is fine.

FAIL if the latency figures are missing or replaced by a qualitative claim ("significantly faster"); if the datastore constraint is missing; or if the ADR points to the Slack thread, a chat, a meeting, "the conversation", "as discussed", or "as agreed" as the place where the reasoning or constraint can be found.
