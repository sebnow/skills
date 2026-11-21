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

## Sentences a reader can follow cold (2026-11)

A brain reply described three parallel mechanisms in one sentence: it
said they "share one shape" without saying what the shape was, joined the
cases with semicolons, and dropped "that" from clauses describing a noun.
shared-shape rebuilds it: a scout reports four documents fetched while
the return is produced. Graders: one-per-line (opus judge) and
short-sentences (regex, no sentence over 25 words above the evidence
heading). Run with `--disallow-tools Agent` and `max_turns: 20`.

Kept, in "How you report": in every paragraph above the evidence, more
than two things go in bullets, not a sentence; one idea per sentence,
under 25 words, split where "and", "so", ";" or ":" joins two, "the
opening summary too".

Model. Everything above this section, and the first rounds below, ran
the brain on claude-fable-5-1. The final rounds ran it on Opus
(`--model claude-opus-5-5`, with `CLAUDE_CODE_EFFORT_LEVEL=high` in the
environment; the trace confirms the model, the effort is unconfirmed),
since Fable is not available on every plan. Compare numbers only within
one model.

On Fable, 3 runs each (one-per-line, short-sentences):

- before (`rev:txprmkru`): 0, 0 of 3.
- rules for decision items only: 0, 0 of 3. Replies put the four cases
  in a context paragraph, as "The four are A, B, C and D."
- rules for every paragraph, cases as bullets, joined sentences split:
  3, 2 of 3.

On Opus, 6 runs each:

- before: 2, 0 of 6.
- A, the rules for every paragraph: 6, 4 of 6.
- B, A plus "the general statement first, then its cases" and "say in
  one sentence what they have in common": 5, 4 of 6. A kept.
- A plus "the opening summary too" (kept): 6, 5 of 6. Most misses on A
  were the lead sentence joining two ideas with ":" or "so".

Dropped for no measured effect: stating what the cases share before
them (the fault never reproduced; B above did no better); a screen
budget counting items rather than words; "the verb follows its
subject"; "active voice"; "keep 'that'", with or without a Not/But for
it (its grader did not move, so it went too).

Regressions on the final wording, Opus, 3 runs: buried-ask 9, 9, 9 of 9;
label-only 2, 4, 3 of 4 (named-rules-stated 2 of 3). label-only on
master's brain.md on Opus gave 2, 2, 4 of 4, so the misses are not the
change.

## Routing work to workers (2026-11)

A Codex brain session made nine shell calls of its own before its one
scout dispatch: it inspected configuration, ran evals, wrote fixtures and
investigated dispatch errors, then did git and release chores. The user
had to say "You're the brain. Dispatch to subagents." Earlier smoke
tests missed it because each prompt named its worker's job and the
grader asked only whether that worker started.

The route-* cases rebuild it on tally, a small Python CLI with ADRs and a
git commit. Prompts are what a user types, with no role named. Bash,
Edit and Write are allowed, so the brain can do the work itself.
`tools/root-work.py` counts the root's own tool calls from the kept
trace. A run's workload is ok when the root makes no Edit or Write, at
most two Read, Grep and Glob calls together, and no Bash in
route-question or at most one elsewhere. The graders check only which
worker ran: a `tool_used` grader counts the workers' calls too (a
senior's Edit failed a root "no Edit" grader).

- route-question: a factual question, for the scout.
- route-settled: build ADR 3 (a CSV export), whose decisions are all
  recorded but no edit is spelled out.
- route-release: two version bumps, a changelog heading, the tests and a
  commit, every edit spelled out, for the junior.
- route-midnight: a fix whose behaviour the user hands to the brain to
  decide.
- route-investigate: a report nobody can reproduce, for the senior.
- route-tiny: a typo in a named file, which the brain should fix itself.

Kept, in "What you do yourself" and "What you hand off" of brain.md:
hand inspection, searches, commands, test runs, edits and git chores to
the cheapest capable worker, short chores too; dispatch before you
inspect; pick the worker by who still has to decide something (the
user's wording, with a new `mid` worker between senior and junior); do
not re-read what a worker cites; do not diagnose a worker's environment
blocker yourself. The junior moved to Haiku.

claude-fable-5-1, opus judge, 3 runs per cell. Workload ok, then the
workers each run dispatched:

| case | master | first draft | kept |
| --- | --- | --- | --- |
| question | 0/3; none | 0/3; scout x3 | 3/3; scout x3 |
| settled | 0/3; junior, senior x2 | 0/3; mid x3 | 2/3; mid x3 |
| release | 0/3; none | 0/3; junior x2, mid | 3/3; junior x3 |
| midnight | 0/3; none x2, junior | 1/3; mid x3 | 3/3; mid x3 |
| tiny | itself 3/3 | itself 3/3 | itself 3/3 |

On master the brain did the work itself in 8 of the 12 runs that call
for a worker. The question
runs made 10 to 13 calls each; the release runs made every edit and ran
the tests and git. The first draft had the routing rule and the mid but
not the last two rules. Its scout runs re-read the four files the scout
cited. Its release runs ran two or three git commands at the root,
working out why git failed in the sandbox. Both stopped with the two
rules added. The one kept settled miss made three framing calls, one
over the limit. All three kept settled runs built an export that passed
its tests and matched ADR 3 (checked by hand: header, order, quoting,
rounding).

The midnight brain never sent the senior. It decided the day the time
counts toward, wrote the decision into the brief, and sent the mid. The
user had handed it that choice, so this follows the routing text.
route-investigate was added to reach the senior and ran only on the kept
text, 3 runs. The brain sent the senior first in 3 of 3. Workload was ok
in 1 of 3: two runs listed the repository with two shell calls before
dispatch. In one run the brain then decided by itself that invoices round
to the quarter hour, a billing change nobody had decided. It made the
edit and ran the tests itself, as "cheaper to do than to brief".

The release case cannot commit: git is not executable in the eval
sandbox. The case still shows who attempts the commit and what the brain
does with the blocker.

Regressions on the kept text, Fable, graders passed per run, against
master's brain.md in an evals-tmp arm (mid.md was on disk for both):

- settled-ask, of 3: kept 3, 2, 2, 3, 3; master 3, 3, 2. The kept misses
  were conventions-not-asked, which "Asking less" has at 3 of 3 on both
  arms. One miss dispatched a mid at once and asked about the ledger ADR
  without reading the conventions; the other read them and asked which
  ADR supersedes it.
- label-only, of 4: kept 3, 0, 2, 4, 0; master 2, 0. One kept 0 sent a
  junior an edit nobody asked for and ended on a short note about it,
  which the last_message graders then read. The other 0s, on both arms,
  are the usual miss.
- buried-ask, of 9: 9, 7, 9.
- answer-fidelity, of 4: 4, 4, 4, with no-worker-dispatch now counting
  the mid.

A five-hour usage limit ended two runs; they are left out above.

Junior on Haiku (claude-haiku-5-5) against Sonnet, junior cases, 3 runs
each, runs at full score: gap 3/3 and 3/3; report-shape 3/3 and 3/3;
scope 3/3 and 3/3; verify 2/3 and 2/3 (both miss output-in-report
once); stale-line 3/3 on Sonnet, 8/9 on Haiku over 9 runs.

## Watch

- Since the routing change the brain dispatches more readily. In
  settled-ask it once briefed a mid before reading the conventions it
  then asked about. In label-only it sent an edit nobody asked for. In
  route-investigate it decided a billing rule itself. If these recur,
  the routing text needs a line that dispatch waits for the decision.
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
