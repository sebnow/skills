# claude plugin eval, the parts a skill author needs

Full reference: https://code.claude.com/docs/en/plugin-evals.md. This page
holds the shape of a case and the flags this plugin's suites use, so you do
not fetch the docs for every case. `claude plugin eval --help` has the flags.

## Layout

```
<plugin>/evals/<skill>/<case>/
  prompt.md            # frontmatter + the user prompt; or case.yaml below
  case.yaml            # needed for context.* (fixtures, resumed history)
  fixture.sh           # scaffold script, copies resources/ into the sandbox
  resources/           # files the agent finds in its working directory
  history.jsonl        # transcript to resume, so the prompt is the next turn
  graders/*.md         # one file per check
```

Worked examples, relative to the plugin root:
`evals/review-writing/vague-pointer` (prompt.md only),
`evals/review-code/generated-embed-workaround` (case.yaml, fixture, trace
grader), `evals/review-writing/leakage-origin` (history_file).

## prompt.md frontmatter

`name`, `tags`, `runs` (default 3), `max_turns`, `timeout_seconds`,
`allowed_tools`. Without `--allow-tools` at run time, only read-only tools
are allowed: Read, Glob, Grep, Skill, Agent, TodoWrite and the Task tools.
Bash, Write, Edit and web tools need `--allow-tools`; prefer asking for the
artifact in the reply.

## case.yaml

```yaml
schema_version: "1.1"
name: <case>
tags: [<skill>, <failure>]
runs: 3
execution:
  max_turns: 15
  timeout_seconds: 300
  allowed_tools: [Read, Glob, Grep, Skill]
  prompt: |
    ...
context:
  scaffold_script: fixture.sh    # runs only with --scaffold
  history_file: history.jsonl
```

## Graders (frontmatter in graders/<name>.md, rubric in the body)

- `type: llm`, `focus: last_message | trace | files | {source: file, path}`.
  Three judge votes, majority wins. Body is the rubric: one paragraph of
  context, then a PASS condition and a FAIL condition, each concrete.
- `type: regex`, `target: last_message | trace | files | {source: file,
  path}`, `pattern`, `flags`, `match: contains | not_contains | count:N`.
  It needs no judge calls and gives the same verdict on every run, so use it
  for long output and for checks that something is absent.
- `type: tool_used`, `tool`, `input_match` (regex on the JSON input), `min`,
  `max`. `min: 0, max: 0` asserts the tool was not called.
- `type: tool_order`, `before`, `after`. `type: file_exists`, `path`.
- `weight` (default 1). `arm: with-only` keeps a grader out of the without
  arm's score; `arm: both` forces it in. A `tool_used: Skill` grader is
  treated as with-only automatically and reported as "plugin fired", which is
  how you tell a zero delta from a skill that never loaded.

Skill-fired indicator, one per case:

```yaml
---
type: tool_used
tool: Skill
input_match: '"skill"\s*:\s*"(?:[\w-]+:)?<skill-name>"'
---
```

## Running

```sh
claude plugin eval . --tag <skill> --scaffold --judge-model sonnet -j 3 \
  --no-publish --max-cost-usd 10 --json /tmp/<skill>.json
```

Ablation is `with-without` by default, so each case reports WITH, W/OUT and
delta. `--ablation none` runs one arm; use it for the first baseline of a
skill that does not exist yet. `--runs 10` when a 3-run result decides a
deletion. `--model` pins the model under test, so a change in the default model does
not look like a regression in the skill. `--case <glob>` filters by name. Results land in
`evals/results/<timestamp>/aggregate-result.json` (ignored by VCS) plus
`report.html`; the JSON has `cases[].aggregates.{score,scoreWithout,delta}`
and per-run grader verdicts with the judge's explanation.

A run's score is the weighted fraction of graders that passed; a case's
score is the mean over runs; the exit code is 1 when any case's with-arm
score is under `--threshold` (default 1.0). Delta never affects the exit
code, so read it yourself.
