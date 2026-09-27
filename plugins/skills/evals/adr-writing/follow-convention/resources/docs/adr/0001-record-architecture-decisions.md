# Record architecture decisions

## Status

Accepted

## Context

Decisions about ledger's storage, APIs and deployment are made in design
reviews and pull requests, and the reasoning is hard to find a year later.
New team members ask why things are the way they are, and the answer depends
on who is still around.

## Decision

We record significant architecture decisions as ADRs in `docs/adr/`, one
Markdown file per decision, numbered sequentially as `NNNN-short-title.md`.
Each ADR has Status, Context, Decision and Consequences sections. An accepted
ADR is not rewritten; a later decision that changes it gets a new ADR that
supersedes it.

## Consequences

Reviewers can ask for an ADR when a pull request makes a decision that is
hard to reverse. The ADR directory becomes the place to look before changing
storage, API contracts or deployment topology.
