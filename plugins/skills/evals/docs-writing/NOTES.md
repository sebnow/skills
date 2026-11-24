# Repository docs evals: notes

## What the cases are for

The user wants a skill that documents the repository it runs in, based on a
system-documentation skill from a central docs repository. It should follow
Diátaxis, write clearly, use diagrams, and stay at the level of the
architecture instead of repeating what the code says.

Failures the cases look for:

- Pages that restate the code: source paths, class and function names,
  configuration tables, constants.
- Pages that mix Diátaxis types, or serve only one of the two default
  audiences.
- Text diagrams instead of Mermaid, and lifecycles without a state diagram.
- Words and marks the style rules ban.
- Pages outside the repository's documentation convention.

Cases:

| case | fixture | prompt |
|---|---|---|
| document-whole-repo | inkwell, README only | "write proper documentation for this repo" |
| docs-site-convention | inkwell, plus an MkDocs site in `documentation/` with a style guide that asks for Title Case headings | "add docs on how inkwell's conversion jobs and webhooks work" |

Both fixtures copy inkwell from `../architecture-visualization/resources`.
To run these cases in another plugin, copy that directory too.

Graders:

- Most checks are `tool_used` on Write with an `input_match` on the
  JSON-encoded input. They see every page the agent writes, and they never
  see the skill body. A `max: 0` grader needs `min: 0` too, because `min`
  defaults to 1.
- `nav-updated` is a regex on the trace, because the agent writes
  `mkdocs.yml` with either Edit or Write.
- `high-level` and `diataxis` are opus judges on the trace. The judge sees
  only the first and last 12 messages, so it can miss pages written in the
  middle of a long session.
- `title-case-headings` checks that the repository's style guide wins over
  the skill's sentence case.

## Runs (2026-11-24, claude-opus-5-5, judge claude-opus-5-5)

Bash is left out of `allowed_tools`. With Bash granted, every run failed to
start: the sandbox could not exclude a credentials file that an AWS
`credential_process` points at. The skill's shell checks therefore do not
run here; the with-skill runs made them with the Grep tool instead.

### Without the skill

The first baseline ran from an empty plugin with `--ablation none`; the
Results table below is from the later with-without runs.

document-whole-repo, 3 of 3 runs:

- Pages went to `docs/` (architecture, api, configuration, operations,
  development), with every Diátaxis type mixed into each.
- An ASCII box diagram of the parts, and not one Mermaid block.
- Source paths (`inkwell/config.py`), class names, and a table of every
  environment variable with its default.
- Em-dashes in tables; both judges voted FAIL in every run.

docs-site-convention, 3 of 3 runs: pages went to `documentation/`, the nav
was updated, and the headings followed Title Case. No run drew the job
lifecycle as a state diagram. Five of the six baseline runs across both
attempts named source files or classes.

### With the skill

document-whole-repo scored 1.0 in 3 of 3 runs. The skill fired in every
run; two runs also loaded architecture-visualization. Each run wrote an
overview, an architecture page with a parts flowchart, a glossary,
lifecycle and flow explanations, how-to guides, an API reference, a feature
map, and an entry page.

docs-site-convention scored 1.0 in 3 of 3 runs. Each run kept the pages in
`documentation/`, wrote `mkdocs.yml` with the new pages in the nav, used
Title Case headings, and drew the job lifecycle as a state diagram.

### Results

From the with-without runs (`claude plugin eval . --case <case>`, 3 runs per
arm):

| case | with | without | delta |
|---|---|---|---|
| document-whole-repo | 1.00 | 0.19 | +0.81 |
| docs-site-convention | 1.00 | 0.76 | +0.24 |

| grader | without (of 3) | with (of 3) |
|---|---|---|
| document-whole-repo: docs-dir | 3 | 3 |
| document-whole-repo: diataxis | 1 | 3 |
| document-whole-repo: high-level | 0 | 3 |
| document-whole-repo: has-diagram | 0 | 3 |
| document-whole-repo: no-ascii-diagram | 0 | 3 |
| document-whole-repo: no-banned-words | 0 | 3 |
| document-whole-repo: no-code-identifiers | 0 | 3 |
| docs-site-convention: in-docs-dir, no-default-docs-dir, nav-updated, title-case-headings | 3 | 3 |
| docs-site-convention: lifecycle-state-diagram | 0 | 3 |
| docs-site-convention: no-code-identifiers | 1 | 3 |

The model already follows a repository's documentation location, its
navigation file, and its style guide without the skill. The skill keeps the
short rule on location because it also sets the `docs/` default, and the
style guide precedence because the skill's own style rules would otherwise
conflict with it.

## Watch

- With the skill, the final reply lists many suspected defects in the code,
  because the skill says to report each one. That is intended here, but it
  makes the reply long.
- The no-code-identifiers pattern names inkwell's classes. Update it if the
  fixture changes.

## Running

From `plugins/skills`, with the plugin loaded (with-without arms):

    claude plugin eval . --case document-whole-repo --scaffold \
      --allow-tools Write Edit --model claude-opus-5-5 \
      --judge-model claude-opus-5-5 --no-publish --trust-plugin -j 3 \
      --keep-temp --json /tmp/docs-writing.json

Run `--case docs-site-convention` the same way. The case names are unique
across the suite, so `--case` selects exactly one. A run of
document-whole-repo costs about 11 USD for both arms.
