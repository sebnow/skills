---
name: review-code
description: "Reviews a code change and reports findings. Use before surfacing or committing a change, on a diff, commit range, or set of files, and when asked to review code. Triggers: 'review code', 'review this change', 'review my diff', 'is this ready to commit', pre-merge review. Not for style-only or formatting nitpicks."
disable-model-invocation: true
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

## Tests that restate the code

Tests are part of the change and must earn their keep. A test earns it by
failing when a decision the code makes goes wrong. A test that restates the
implementation can find nothing the code does not already say; it fails only
on intended changes, and so creates churn.

For each test, name the decision it checks, where that decision is made, and
which other test in the change already fails for the same break at no more
cost. It is suspect when:

- what it asserts is fixed by the code's own text — the order of elements a
  template lays out, exact markup or wording, the sequence of calls a
  function makes;
- another test no slower than itself already fails for the same break, and
  this one reaches no edge case of its own. Checking the same values more
  directly, or under another name for the decision, is not an edge case of
  its own.

Order or values that logic decides — a sort, a grouping, a computed total —
are behaviour, not restatement; a test of them stays unless another test
already covers it. A fast test that stands in for most of a slow or paid one
stays beside it.

Report the test and what it restates or which test already covers it. Fix:
delete it, or rewrite it to check a decision the code makes. Pinning the whole
output as a literal or snapshot restates the code again.
