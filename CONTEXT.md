# Session context: automating the code-quality loop

Ephemeral working memory for the `sebnow:brain` session. Not part of the
plugin; purge from history along with `RESEARCH-POTETO.md`. Nothing in the
plugin (skills, agents, workflows, README) or in commit messages may cite
this file, the research note, or pstack.

## Goal

Automate as much as possible of the verification loop for software
development, where "verified" means the result is well written with low
technical debt, not merely that tests and linters pass. The loop must
improve existing code regardless of who wrote it; agents do well on
greenfield work and fail where debt and bad architecture already exist.

## Inputs

- `RESEARCH-POTETO.md`: research into pstack (Lauren Tan). Local checkout
  at `~/code/plugins/pstack`. pstack's workflows are prompt playbooks; its
  only code is a state store (`skills/poteto-mode/scripts/orch/`, verdict
  vocabulary enforced in `store.ts:286`) and a plan linter. It classifies
  intent (`/poteto-mode <goal>` picks a playbook) but does not thread it
  as data. History tidy rules borrowed from `playbooks/opening-a-pr.md`
  (Commits paragraph), `principle-sequence-verifiable-units`, and
  `playbooks/refactoring.md` step 8.
- Claude Code dynamic workflows (docs: code.claude.com/docs/en/workflows.md):
  JavaScript scripts with `agent()`, `pipeline()`, `parallel()`, `phase()`,
  `log()`, `args`, `budget`, `workflow()`; schema-forced JSON output;
  `agentType` selects plugin agents; scripts ship in a plugin's
  `workflows/` dir and run as `/sebnow:<meta.name>`. No filesystem or
  shell from the script; no `import()`; `Date.now()`/`Math.random()`
  throw. Runs count against subscription usage. Installed version 2.1.280;
  workflows are enabled.

## What was built

Commits on master, unpublished, both rewritten in place as the design
changed (no intermediate versions kept):

- `uxkkslol` "workflows: add improve-code review-and-fix loop":
  `plugins/skills/workflows/improve-code.js` (~400 lines) and the lens
  agents `plugins/skills/agents/reviewer-{structural,general,tests,duplication}.md`.
- `uklnwtmy` "workflows: add a writing lens to improve-code":
  `reviewer-writing.md` preloads the `review-writing` skill via
  `skills:` frontmatter (works only for skills without
  `disable-model-invocation: true`) plus a comments-restate-code check;
  reach = scope. Preload untested in a run.
- `ulwuqynu` "readme: document the improve-code workflow": a high-level
  Workflows section in `README.md` (user asked for high level only).

### `/sebnow:improve-code` shape

1. Scan (haiku, low effort): runs the VCS lookup command verbatim, finds
   the check command, reports layout, conventions, docs, `has_tests`.
   Facts only; its brief is prepended to every later prompt.
2. Review: one reviewer per lens over the whole scope (not per file or
   unit; a per-unit fan-out produced 31 Opus reviewers on 9 units and was
   stopped by the user). Each lens is a plugin agent selected via
   `agentType`; lens text lives in the agent file. Reach per lens: tests
   = scope plus code under test; general = one hop (callers, callees,
   owning package); structural and duplication = whole codebase;
   writing = scope (docs, comments, messages, commit messages in a range).
   Findings: `locations` (schema pattern `^[^:\s]+:\d+$`), problem,
   consequence, `evidence` in {traced the path, ran it, reasoned,
   hypothetical}, evidence_detail, severity, fix; plus
   `read_outside_scope`. Tests lens skipped when `has_tests` is false.
3. Triage: one agent, buckets act-on / consider / dismissed with a reason
   each; `recurs` names a previously fixed finding; a recurrence ends the
   run as escalated, never retried.
4. Fix: one fixer per round, findings in order, one commit per finding,
   check must pass or the change is reverted. Fixes at the cause,
   anywhere in the repo.
5. Re-review starts from the files the fixes changed.
6. Tidy (only if a fix landed): absorb-or-separate per fix; jj
   `absorb --from <fix> --into <range>` then abandon emptied commits;
   git `--fixup` + `--autosquash` or `git absorb`; order commits so a
   reviewer can replay them; test and fix stay in one commit when the
   test would be red alone; project commit style via the commit skill;
   every rewritten commit re-checked (`git rebase -x`; jj `new`/check/
   `abandon`); failures reported, never fixed; never edits code.

Stop conditions: no findings, no act-on, recurrence, no fix landed,
`rounds` (default 3). Args: `paths` or `revs` (range reviewed as one set,
commits matter only as absorb targets), `rounds`, `check`, `model`,
`models` and `efforts` per stage (`scan`, `review`, `triage`, `fix`,
`tidy`); only scan has a cheap default. JS identifiers camelCase; JSON
keys snake_case.

## Verification status

- The first shape (unit x lens) ran once end to end on a seeded Go fixture
  via a headless `claude -p` session: 3 seeded problems found, triaged,
  fixed in three commits, `go test` passing. Headless run cost about
  0.60 USD.
- Everything since (reach, duplication lens, scan, per-lens fan-out, single
  fixer, tidy, per-stage models) is verified with `node --check` only.
  The user's one interactive run was on the unit x lens shape and was
  stopped for cost.
- Unknowns: whether `effort` and schema `pattern` work inside a workflow
  agent; whether the copied structural lens text reaches its reviewer
  (cause found and fixed: `disable-model-invocation: true` blocks
  `skills:` preload; the flag is dropped from `review-code` and
  `reviewer-structural` now preloads it instead of carrying a copy;
  untested in a run);
  whether `/sebnow:improve-code` invoked interactively calls the Workflow
  tool reliably (in `claude -p` with a haiku driver it did not).

## Decisions and preferences from the user

- Workflow in code, not brain text or a skill; brain is command and
  control only, and agents are competence/model config, not workflow.
- Rubrics are the hard part and do not exist yet; even a plain
  fresh-context review catches slop, so orchestration first, calibration
  later. Evals wanted, but real fixtures must be collected first,
  incrementally.
- Scope is a starting point, not a boundary; reach is per lens.
- Compression-oriented: duplicated code is a refactor candidate and the
  reviewer must find the other copies.
- Deterministic lookups (VCS) as explicit commands, not inference; cheap
  model for fact gathering, reasoning model for judgment.
- Tidy must re-verify every commit and work on git as well as jj.
- Keep cost low: the account is near its budget limit. No workflow or
  `-p` runs from agents; the user runs interactively. `claude -p` may be
  API-billed rather than subscription (user's statement; unverified).
- README stays high level. No superfluous comments. Commit style
  `scope: lowercase imperative`, not conventional commits.
- Not adopted from pstack: arena/architect (needs model-vendor
  diversity), reflect (skill-writing already starts from a failure),
  TSV decision log, no-comments agent, sticky mode with principle naming.

## Next steps

1. User runs interactively on real code with debt:
   `/sebnow:improve-code {"paths":["<dir>"],"rounds":2,"models":{"review":"sonnet"}}`.
   Expect at most 4 review agents per round. Bring back the returned
   object.
2. Judge from it: structural and duplication lens output, dismissals,
   the re-review round, tidy result, and sonnet-vs-opus for reviewers
   (open decision: default review model).
3. Then: collect real fixtures for evals; consider a requirements lens
   later; decide whether the `review-code` skill text and
   `reviewer-structural` should share one source.
4. Open minor items: fixers report commit hashes not jj change ids; check
   output may show `(cached)` (pass `check` with `-count=1` for Go);
   `sebnow:review-writing` not run on agent text or README.

## Workers

Senior agent from this session (id `a9369ae7562e07cc0`) built the
workflow and holds its context; resumable with SendMessage. Brain did
the scan stage, the per-lens fan-out, and the README cut itself.
