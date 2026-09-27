---
name: driver-stated
tags: [adr-writing, driver]
max_turns: 12
allowed_tools: [Skill, Agent]
runs: 3
---

Need an ADR for orders-service. We're adopting an outbox for the events we send to other services (`order.placed`, `order.cancelled`, `order.refunded`).

How it works: `OrderService.Place()` and friends insert a row into a new `outbox` table (`id`, `aggregate_id`, `event_type`, `payload jsonb`, `created_at`, `published_at`) in the same transaction as the order change, instead of calling `publisher.Publish()` after the commit like they do now. A new `outbox-relay` deployment (same repo, `cmd/outbox-relay`, 1 replica, leader lock via `pg_try_advisory_lock`) polls every 500ms, picks up to 100 unpublished rows ordered by `id`, publishes them to the `orders` exchange on RabbitMQ with publisher confirms, and sets `published_at`. Published rows are deleted after 7 days by a nightly job. Delivery becomes at-least-once; fulfillment and billing already dedupe on `event_id`, so that's fine. Expected extra latency is ~250ms average, the 500ms poll interval worst case.

For context, ops is still grumpy about the 16th, when fulfillment had to hand-replay 1,214 orders after the rolling deploy of orders-service. We've got ~9k orders/day, peaks of ~40/min.

Put the ADR in your reply.
