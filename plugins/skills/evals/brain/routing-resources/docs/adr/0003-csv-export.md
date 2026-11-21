# 3. CSV export of the time log

Status: accepted (2026-11-25)

## Context

Two clients want the raw entries behind each weekly report, to load into
their own spreadsheets.

## Decision

- A new command, `tally export LOG --client NAME`, writes one client's
  entries as CSV to standard output. Without `--client` it writes every
  client's entries.
- Columns, in this order, with a header row: `date`, `client`, `task`,
  `start`, `end`, `minutes`, `billable_minutes`. Dates are ISO 8601
  (`2026-03-02`), times `HH:MM`.
- `minutes` is the raw duration; `billable_minutes` uses the same rounding
  as the weekly report (ADR 2), so the two always agree.
- Rows are sorted by date, then start time.
- The CSV uses the standard library `csv` module with its default dialect.
  A task containing a comma is quoted by the writer, not stripped.
- An unknown client is not an error: the export prints only the header.

## Consequences

The export lives beside the report and invoice modules and gets its own
tests. The README lists the new command.
