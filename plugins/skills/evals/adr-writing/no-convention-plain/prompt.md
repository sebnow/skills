---
name: no-convention-plain
tags: [adr-writing, template]
max_turns: 12
allowed_tools: [Skill, Agent]
runs: 3
---

Starting ADRs in this repo. There are none yet: no docs/adr directory, no template, nothing, so the format is up to you. First one: checkout-api will call pricing-service over gRPC instead of the REST/JSON we use everywhere else.

Background: checkout calls pricing 4-6 times per checkout (quote, promo evaluation, tax estimate). The p50 of those REST calls is 18ms and most of that is JSON encoding and decoding of the line-item payloads, which run up to 300 items for B2B carts. We want one typed contract both teams review in one place (`proto/pricing/v1/pricing.proto`), and deadline propagation so a slow pricing call can't hold a checkout past its 2s budget. Both services are Go. The public API stays REST. The costs are that we lose curl-ability (we'll need grpcurl and server reflection enabled outside prod) and the internal ALB needs HTTP/2 target groups for pricing.

Give me the ADR in a markdown code block so I can paste it, and tell me what filename to save it under.
