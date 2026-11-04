# Brain evals: notes

## Asking less (2026-11)

A brain session (orchestrator project, 2026-10-25) asked the user 39
questions, 14 of them answerable from a file it had already read or the
user's own earlier words, and used names from worker reports as if the
user had seen them. Three fixes were tried:

- Check the session's sources before asking, ask at most two questions
  per reply, and end each decision's context with a `Checked:` line
  (`ustkwyxu`, the "# Before you ask" section of `agents/brain.md`).
  Kept. On settled-ask it passed 8 of 8 runs against 1 of 6 for the
  prompt before it (`rev:onyvwynn`); the last 3 against 3 were checked by
  hand against the first reply. A variant without the `Checked:` line
  did as well on settled-ask but stated rules less often on label-only
  (1 of 8 against 3 of 6), so the line stays.
- Introduce every name taken from a worker report. Dropped with its
  case, term-provenance: the brain passed that case about half the time
  without the rule, and the rule moved it from 5 of 12 to 8 of 12, which
  is within noise. A case worth building would start from a worker's
  draft that coins opaque terms the brain then repeats ("Entry",
  "Harness Line", "key prefixes").
- Send a scout the draft questions before replying. Not written: with
  the first fix, settled-ask has nothing left for it to catch.

Results on master after the first fix, 3 runs each, opus judge:
buried-ask 9, 9, 8 of 9 graders (before: 8, 9, 9); answer-fidelity-option
3 of 3; label-only, named-rules-stated 2 of 3 (before: 1 of 3).
label-only fails often on every arm: a reply that names a rule by subject
("how carried hours count") instead of saying what it requires fails.

## Watch

- Two replies with the first fix decided settled-ask's sold-vessel
  point, which the case counts as open. Nothing grades that; if it
  recurs, add a grader, since "do not ask what is settled" can turn into
  settling open points.
- settled-ask's fixture has no vessel in its ledger, and replies often
  ask about that. Fair, but it adds noise.

## Harness

- Run from `plugins/skills`, after `sync-system-prompt.sh`:
  `claude plugin eval . --case <case> --ablation none --scaffold --runs N
  --judge-model claude-opus-5-5 --no-publish --trust-plugin -j 1
  --keep-temp --json <out>.json`.
- Judge with opus. A sonnet judge failed replies that an opus judge and
  a reader passed. `tools/regrade.sh` regrades kept runs with another
  judge; `tools/probe-grader.sh` tries one grader on saved replies;
  `tools/trace-replies.py` prints what the brain said in a kept trace.
  Probe a new llm grader on known good and bad replies before trusting
  it.
- An llm judge on `focus: trace` sees only the first and last 12
  messages. settled-ask needs `trace` anyway: Agent is allowed, so the
  last message can come after a worker reports back.
- `allowed_tools` does not block Agent; `--disallow-tools Agent` does.
- `-j 2` on one case makes its runs share a transcript; use `-j 1`.
- Arms pinned to other revisions: copy the case to
  `evals-tmp/brain/<name>`, set its `name`, run a copy of
  `sync-system-prompt.sh` there with its own `cases` list, and pass
  `--eval-dir evals-tmp`. Keep `evals-tmp` out of commits.
- The script that generated settled-ask's history was lost; edit
  `history.jsonl` directly.
