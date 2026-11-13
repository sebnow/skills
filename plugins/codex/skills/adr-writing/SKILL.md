---
name: adr-writing
description: "Writes and updates architecture decision records (ADRs). Use when asked to write, create, draft, or update an ADR, or to document or record an architectural decision. Checks whether the decision warrants an ADR, keeps accepted ADRs unchanged, keeps implementation detail out, and names the file to fit the project. Not for design docs, specs, or RFCs."
---

# ADR Writing

An ADR records why a decision was made at a point in time: the constraints
that shaped it and the alternatives it rejected. The code records how the
decision was carried out, so the ADR does not repeat it.

## Which decisions get one

A decision gets an ADR when it is hard to reverse, when changing it later has
a meaningful cost, or when a future reader would ask "why did they do it this
way?". A change that is cheap to reverse and touches one place, such as a
dev-only config value, fails all three; its commit message is enough.

When the request is for such a change, do not write the ADR, not even a short
version. Tell the user it does not need one, say which of the three tests it
fails, and ask whether they still want it. Do this even when the user asked
for the ADR outright: a log full of small records buries the ones that matter.

## Accepted ADRs are not edited

An accepted ADR is the record of what was decided and why at that time.
Editing it rewrites the history a later reader relies on. Two edits are
allowed:

- changing the status line to point at the ADR that supersedes or amends it
  ("Superseded by ADR-0007");
- fixing a typo or a dead link.

Everything else goes in a new ADR: a new rule or convention, a "Related" or
"See also" link, anything that happened after the decision. When asked to add
a rule to an accepted ADR, write a new ADR that amends or supersedes it, and
offer the status-line pointer as the only change to the old one. If the new
content is a registry that changes often (names, owners, endpoints) rather
than a decision, propose a living document instead.

## Implementation detail stays out

Schemas, DDL, config, and code appear only to show the shape of an option: at
most one short snippet per option, and never the artefact that ships. The ADR
is a guideline, and the implementation keeps room to move; a copied schema
goes stale with the next migration. Name the boundary and the
responsibilities, such as which data gets a real column and what makes a
field move out of a JSON blob, instead of listing columns, indexes, or
constraints. Point to the migration or config file by path when the reader
needs the detail.

## Naming the file

If the project already has ADRs, match their directory, filename pattern, and
section layout.

If it has none, name the file `YYYY-MM-DD-<need>.md` with today's date, after
the need the decision answers rather than the solution chosen:
`session-storage`, not `use-redis-for-sessions`. The solution may be
superseded while the need stays. Leave out sequence numbers; two branches that
each add ADR 0005 collide when they merge.

## Revisit if

The ADR may end with one line naming what would make the decision worth
reopening, such as "Revisit if peak writes pass 5,000 per second." Keep it to
that line, not a section.
