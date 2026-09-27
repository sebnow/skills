# PostgreSQL for ledger storage

## Status

Accepted

## Context

Ledger stores double-entry postings for about 3,000 merchant accounts. Every
transfer writes at least two postings, and the sum of postings for a transfer
must be zero at every point a reader can observe. Balances are read far more
often than they are written: roughly 40 reads per write at peak.

The team runs PostgreSQL for two other services and has on-call experience
with it. DynamoDB was the other candidate; it offers transactions across up to
100 items, but balance queries over a date range would need a second index
per access pattern.

## Decision

Ledger stores postings and accounts in PostgreSQL 15 on RDS. A transfer and
its postings are written in one transaction with `SERIALIZABLE` isolation.
Balances are computed from postings, with a nightly snapshot table so that a
balance read scans at most one day of postings.

## Consequences

Transfers are atomic without application-level compensation. Write
throughput is bounded by a single primary; at the current 60 writes per
second this is not close to a limit. Serialization failures must be retried
by the caller, and the API layer does this up to three times.
