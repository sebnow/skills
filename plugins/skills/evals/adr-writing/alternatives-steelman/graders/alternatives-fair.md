---
type: llm
focus: last_message
weight: 1
---
The conversation before this ADR compared two options for dispatching invoice PDF render jobs in depth: Postgres LISTEN/NOTIFY with a jobs table, which was chosen, and a Redis stream, which was rejected. The Redis stream's real strengths were its consumer groups, with the pending entries list and XAUTOCLAIM giving recovery from crashed workers without extra code, and its throughput headroom. It lost for concrete reasons. Enqueueing to it cannot be atomic with the Postgres transaction that creates the invoice, so it needs an outbox table, which is essentially the jobs table anyway. The existing Redis is an allkeys-lru session cache that could evict entries, so it would also need a new non-evicting instance. The user mentioned Kafka only in passing ("never looked at it"); it was never evaluated.

Judge only the ADR in the reply.

PASS if all of these hold: the ADR presents the Redis stream as an alternative, states at least one of its genuine strengths (consumer-group delivery tracking or crash recovery via pending entries or XAUTOCLAIM, or throughput headroom), and gives at least one of the concrete reasons it lost (the dual write or outbox needed for transactional enqueue, or the eviction-prone cache requiring a new instance); the chosen Postgres option appears with its costs or limits (for example non-durable notifications needing a poll fallback, or LISTEN needing direct connections outside PgBouncer); and Kafka is either absent or explicitly marked as not evaluated.

FAIL if Kafka is listed or discussed as an alternative that was considered or rejected (with or without a reason); if the Redis stream is rejected with only a vague reason such as "too complex", "more operational overhead", or "not a good fit" and no concrete mechanism; or if the Redis stream appears with no strength stated.
