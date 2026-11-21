# 2. Weekly reports round each entry to the quarter hour

Status: accepted

## Context

Clients agreed to be billed in quarter hours, and the weekly report is what
we send them before the invoice.

## Decision

The weekly report rounds each entry's duration up to the next quarter hour
before summing. Rounding lives in `tally/rounding.py`.

## Consequences

A ten-minute call shows as fifteen minutes.
