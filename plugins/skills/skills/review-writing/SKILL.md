---
name: review-writing
description: "Reviews a written document for LLM slop and context leakage before it is treated as finished. Use as the final step of writing any standalone artifact - ADR, README, design note, spec, PR description, commit body, or long code comment - and when asked to review a document or check whether it stands on its own. Reviews the document from a reader's position, in a fresh context that never saw the conversation, then revises. Catches cliches, em-dashes, empty intensifiers, overclaims, defensive negatives, and sentences that answer the conversation rather than a reader. Not for chat replies, summaries, or scratch notes."
---

# Review Writing

A document is read without the conversation that produced it. The writer, still
holding that conversation, cannot see what leaked from it. The fix is to review
from the reader's position - a context that never saw the discussion.

## How to review: a fresh context

Do not review the draft yourself. You still hold the conversation, so you will
read past the leaks. Instead:

1. Launch a subagent (the Agent tool, a fresh general-purpose agent - not a
   fork, which would inherit this context).
2. Pass it only the document text and the instructions below. Do not pass the
   conversation, the user's prompt, or your notes.
3. Apply its findings to the draft, then present the revised document.

## What the reviewer flags

Give the subagent this list. For each hit it quotes the sentence and names the
problem.

- **Context leakage.** Any sentence that answers, defends against, or narrates
  something not present in the document itself: "not speculative", "this is not
  an X framework", "sources that do not exist", "we considered ... and
  rejected", an intensifier like "no matter what" or "by construction" that
  implies an objection the document never raises.
- **Overclaim.** "structural", "by construction", "guaranteed", "impossible",
  "cannot" for a rule a check or convention enforces. Propose the accurate verb
  ("enforced at the write path").
- **Circumstance and history.** Project status ("personal tool", "prototype",
  "for now") or narrative ("after discussion", "working through a prototype")
  that a reader in two years cannot verify or will find false.
- **Stilted negation.** An absence compressed into "X has no Y" or "A and B
  share no Z", or a riddle that withholds its subject until after a colon
  ("has to generate what the model does not store: the offsetting posting").
  Propose the plain form a person would say: "cash does not have lots", "the
  two ends do not share an identifier", "the model does not store the
  offsetting posting, so the export generates it".
- **Dead pointer.** A reference the reader cannot follow: "documented
  elsewhere", "described in a separate document", "see the relevant section".
  Name the location (URL, file, section). If it is not known, ask the user for
  it; do not drop or reword the pointer silently.
- **Dismissed alternative.** A sentence that raises options only to say they
  are unused: "the format accepts a few other characters, which the project
  does not use". When the document already says what is used, delete the
  sentence rather than rewording it.
- **LLM-isms.** Cliches, em-dashes, empty intensifiers, "it reads", "and
  nothing else", padded phrasing that says nothing.

## Apply and repeat

Apply every finding. Then present the revised document. If the reviewer
reported nothing, say so and present as is.
