# 3. Service naming

Date: 2025-07-14

## Status

Accepted

## Context

Workload names in the cluster are inconsistent: `ledger`, `payments-ledger-svc`,
`LedgerAPI`, `ledger-v2`. Dashboards, alerts and on-call runbooks key on the
workload name, and people guess wrong under pressure.

## Decision

Every Deployment, StatefulSet and CronJob is named `<team>-<domain>-<role>`:

- `team` is the namespace short code from ADR-0002.
- `domain` is the business area, one lowercase word (`ledger`, `refunds`,
  `search`).
- `role` is one of `api`, `worker`, `cron`, `gateway`.

Names are lowercase, hyphen-separated, and at most 40 characters. Examples:
`pay-ledger-api`, `pay-refunds-worker`, `cat-search-api`.

## Consequences

Alert rules and dashboards can be templated on the name. Existing workloads
are renamed when they next have a breaking deploy; the `name-lint` CI check
warns on old names and fails on new ones.
