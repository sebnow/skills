---
name: delegate
description: "Hands long-running or token-heavy work to a worker agent with a self-contained brief. Use when asked to delegate, hand off, or get a subagent to do something, and whenever a task would burn many tokens or run long: multi-file changes, repo-wide sweeps, research, review of a large diff, root-cause investigations. Picks the worker (senior, junior, scout) from how decided the task is and writes the brief for that worker. Do NOT use for a small edit or a lookup with a known location; do those directly."
---

# Delegate

Delegating buys a small context for the caller. It costs a brief, a worker's
run, and a report to read. Hand off only work where the saving is worth
that cost.

## Whether to delegate

Delegate when the work will burn many tokens or run long: changes across
several files, a sweep of the repository, research, review of a large diff,
an investigation with the cause unknown.

Do it yourself when writing the brief would cost more than the work: a
rename, an edit to one file, a lookup with a known location. When asked to
delegate something that small, say it is too small to hand off, and do it
directly.

## Which worker

Pick by how much is already decided, the way you would pick between a
senior and a junior engineer.

| Worker | Use when | What it gets |
| --- | --- | --- |
| `sebnow:junior` | Every decision is made: files, signatures, behaviour, error text, tests, verification | A literal spec. It follows it and stops when the spec is silent |
| `sebnow:senior` | Judgment calls remain: an investigation, a refactor, a feature with unsettled details | A goal and hard constraints. It decides the rest and reports its decisions |
| `sebnow:scout` | The task is read-only and factual: sweep the repo, research a question on the web, mine a transcript for what was decided or what failed | A question and its sources. Several can run in parallel |

If the spec would be longer than the work, use the senior. If you are
choosing the junior, you must be able to write acceptance criteria; if you
cannot, the decisions are not made yet.

## The brief

Build it from the decisions already made. Pass pointers, not content: paths,
symbols, commit ids. The worker reads for itself. Paste something only when
it is not on disk.

Every brief carries:

- The goal, in one or two sentences.
- What is decided. For the junior, everything, written as acceptance
  criteria. For the senior, the hard constraints only.
- The verification command that must pass before it reports.
- The stop rule: stop and report when the brief does not cover something
  you need, instead of inferring.
- The report shape: Changed, Decisions, Verification with the command
  output, Open questions.

A scout brief is the question, the sources (repo scope, URLs, or a
transcript path), and the answer format: a pointer (`path:line` or URL)
with a one-line note.

## After it returns

Read the report. A verification section must show the command and its
output. Treat "tests should pass" as not run, and send the worker back for
the output rather than running it yourself. If the worker stopped on a
gap, answer it and resume the same agent with SendMessage rather than
writing a new brief.
