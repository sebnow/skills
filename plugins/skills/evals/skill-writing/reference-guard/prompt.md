---
name: reference-guard
tags: [skill-writing, references]
max_turns: 12
allowed_tools: [Read, Glob, Grep, Skill]
runs: 3
---

Our http-client skill is getting long. Move the rate limit section out into a
reference file under references/ and keep SKILL.md to the everyday rules. Put
both files in your reply.

```markdown
---
name: http-client
description: "Conventions for calling partner HTTP APIs from Go services. Use when writing or changing code that calls an external HTTP API."
---

# HTTP client

## Timeouts

Every outbound call takes a context with a deadline. Default to 5 seconds for
reads and 15 seconds for writes; partner SLAs are listed in `docs/partners.md`.

## Idempotency keys

Writes that a retry could duplicate (orders, payouts) send an
`Idempotency-Key` header. Derive the key from the domain identifier of the
operation, not from a random UUID, so a retried request after a crash carries
the same key.

## Rate limits

Partners answer 429 with `Retry-After` in seconds, except Northwind, which
sends `X-RateLimit-Reset` as a Unix timestamp and omits `Retry-After`.
Acme's 429 body is JSON with a `retry_in_ms` field that is more precise than
the header; prefer it when present.

Back off with full jitter between the header value and twice the header
value. Cap at 60 seconds. After three consecutive 429s, open the circuit for
that partner for two minutes and fail fast; the `internal/breaker` package has
the state machine, do not implement another.

Batch endpoints count each item against the quota, not each request.
Northwind's batch limit is 500 items per minute, Acme's is 2000 requests per
minute regardless of batch size. The quota is per API key, and staging and
production share a key for Northwind, so load tests in staging can exhaust
production.

Log each 429 with the partner name, the endpoint, and the wait chosen, at
warn level. Do not log the response body; Acme's includes the API key.

## Errors

Map partner errors to the sentinels in `internal/partner/errors.go`.
`ErrTransient` for 5xx and 429, `ErrInvalidInput` for 4xx other than 429.
```
