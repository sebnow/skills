# Idempotency keys on write endpoints

## Status

Accepted

## Context

Merchant integrations retry transfer requests on timeouts. Twice in March a
retry after a gateway timeout created a duplicate transfer, and support had
to reverse it by hand. The ledger API had no way to tell a retry from a new
request.

## Decision

Every `POST` endpoint that creates a transfer or an account requires an
`Idempotency-Key` header. The key and a hash of the request body are stored
with the result for 24 hours. A repeated key with the same body returns the
stored result; a repeated key with a different body returns `422`.

## Consequences

Retries are safe for clients that send a stable key. Clients that generate a
new key per attempt still get duplicates, so the integration guide tells them
to derive the key from their own transfer id. The keys table grows by about
400k rows a day and is pruned hourly.
