---
type: llm
focus: last_message
weight: 1
---
The user asked for an ADR adopting a transactional outbox in orders-service. The prompt mentions, as an aside, that on the 16th fulfillment had to hand-replay 1,214 orders after a rolling deploy of orders-service. Today the service publishes to RabbitMQ after the database commit, so a pod stopped between commit and publish loses the event. That incident is what forced the decision now. An ADR that only describes the outbox pattern leaves a later reader unable to tell why the team took on the relay and the at-least-once delivery at this point.

Judge only the ADR in the reply.

PASS if the ADR states the concrete trigger: events (or orders reaching fulfillment) were lost during a deploy because publishing happened after the commit, with the incident on the 16th or the 1,214 replayed orders or an equivalent concrete description of the loss.

FAIL if the ADR explains the outbox only in general terms (for example "publish-after-commit is a dual write that can lose events" or "to improve reliability") without saying that events were actually lost in a deploy, or if it mentions the incident only outside the ADR.
