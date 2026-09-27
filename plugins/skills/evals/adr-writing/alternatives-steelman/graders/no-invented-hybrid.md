---
type: llm
focus: last_message
weight: 1
---
The conversation before this ADR argued two options for dispatching invoice PDF render jobs: Postgres LISTEN/NOTIFY with a jobs table (chosen, with SKIP LOCKED claiming and a 30s poll fallback) and a Redis stream (rejected). The Redis option, as the conversation described it, needs an outbox table in Postgres relayed into the stream and a new non-evicting Redis instance; describing it that way is part of the option, not a hybrid. No combination of the two options was ever raised, such as Postgres as the queue with Redis as the wake-up signal, NOTIFY feeding a Redis stream, or a jobs table with Redis consumer groups on top.

Judge only the ADR in the reply.

PASS if the ADR presents no hybrid or combination of the Postgres and Redis options as an alternative that was considered, evaluated, or rejected.

FAIL if the ADR presents such a hybrid or combination as a considered or rejected alternative.
