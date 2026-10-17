---
name: review-design
description: "Reviews the design of a component, a set of changes, or an area of an application and proposes restructurings that would improve it without changing what it does, each with alternative designs and their trade-offs. It reports and changes no files. Use when the user asks how to restructure, redesign, or reorganise code (\"how can we restructure the articles package?\", \"how do we remove this duplication?\"), asks whether duplication, coupling, or ownership in some code should change, asks to find architectural or design issues in an application, or asks for a design or architecture review. Not for carrying out a refactor that has been decided, and not a code review for bugs, style, or readiness to merge (use review-code)."
argument-hint: "[what to review] [what you noticed]"
---

# Review Design

The review takes a target: a component, a set of changes, or an area of an
application such as a feature. A transformation is a change to the code's
structure that leaves what the code does unchanged. The report lists
transformations, each with at least two designs and their trade-offs, and lists
observations such as defects and oddities separately from them. The report is
the whole result: the review changes no files.

## The request

The request is in $ARGUMENTS, or in the user's message when there are none. It
names a target, carries a suggestion, or both. A suggestion, such as "these two
look alike" or "this belongs in X", is one candidate transformation: weigh it
like any other, and turn it down with the reason when it is not worth doing.
Without a target, review the changes in progress (uncommitted changes and the
commits on the current branch), or the repository's main components when there
are none, and say what you reviewed.

## Where to look

The target's own code shows only part of its design. Read what calls it and
what it calls: a decision its callers repeat, or one it makes on their behalf,
shows up only there. Name the domain decisions the target encodes, such as
rules, limits, states, and who owns which data, and where each one is made.
Read the context outside the code that is available, such as documentation,
decision records, configuration, and the contracts of other services, for
decisions the code should agree with.

## Alternatives for each transformation

For each transformation, give at least two designs that achieve it, and for
each design what it gains and what it costs. Leaving the code as it is counts
as a design only when its costs and gains are stated too.

## Keep what the code does

A transformation keeps everything a caller can observe: which fields appear in
which state, ordering, errors, and status codes. What the code does for its
users stays the same, including what each response contains. Before reporting a
transformation, read the code that produces each observable output and check,
case by case, that each design produces the same output. Check against that
code, not against your summary of it.

A defect you find is not a transformation. Report it separately as an
observation; fixing it is outside this review.

## Only transform a current problem

Propose a transformation only to remove a problem the code has now.

If nothing qualifies, say so and propose no transformations.
