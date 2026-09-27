---
status: accepted
date: 2026-08-02
---

# Keep every input in one SQLite ledger

## Context

tidewright reads voyage logs (CSV exports from the chartplotter, NMEA
sentence dumps) and turns them into the CCVR operating return. The return
depends on the imported files and on what the operator states about the
vessel. Both have to survive between runs, be removable when an import turns
out to be wrong, and be explainable in the return (ADR-0006).

## Outcome

The ledger is a single SQLite database file. We use the pure-Go driver
`modernc.org/sqlite` to keep builds cgo-free (ADR-0002). Imported files go in
as blobs and are never modified (ADR-0003). Nothing derived is written.

The schema changes only through a new ADR or an amendment to this one, so a
ledger written by one version opens in the next without a migration step.

## Consequences

- One file to back up and to hand to an accountant or surveyor.
- Removing an import is a delete plus a cascade; nothing derived goes stale.
- Adding a table needs an ADR amendment. Adding a column to an existing
  table does too.
