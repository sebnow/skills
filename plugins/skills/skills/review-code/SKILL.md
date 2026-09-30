---
name: review-code
description: "Reviews a code change and reports findings. Use before surfacing or committing a change, on a diff, commit range, or set of files, and when asked to review code. Triggers: 'review code', 'review this change', 'review my diff', 'is this ready to commit', pre-merge review. Not for style-only or formatting nitpicks."
---

# Review Code

## Machinery that treats a symptom

Before judging whether a piece of code works, ask why it exists. Code that
exists to route around the change's own structure is a symptom; the fix is to
the structure, not to the code.

For each generator, generated or copied file, adapter, glue abstraction,
optional dependency, or sync step, name what it produces and who consumes it.
It is suspect when:

- a committed file is derived entirely from other files in the repository —
  a second copy of a source of truth;
- a path climbs out of its component to reach files another component owns;
- it has one producer and one consumer, both in the repository;
- your own fix is a check that it stays correct — that two copies agree, that
  the output is up to date, that a dependency was wired.

Then ask what the consumer would do if it could reach the thing directly, and
what structural change makes that possible — usually letting the owner of the
data or files expose them, or reading the value from a load that already has
it. Report the structural cause. Fix: restructure so the direct mechanism
works, and delete the machinery. Hardening the machinery keeps the cause.
