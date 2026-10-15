---
status: accepted
date: 2026-09-11
---

# Dismissing finished tasks

## Context

The task-lifetime record (0007) kept a task's runner-side state until the
task was dismissed, and
made dismissal the only act that ends a task for good: the runner
deletes the task's record, spool and workspace, and the task takes no
more follow-ups. That record spoke of dismissal only for stopped and
failed tasks ("Until dismissal, a stopped or failed task's workspace
and record are kept"), and the server refused to dismiss any other.

Finished tasks pile up on runners the same way: a task that finished
yesterday holds a workspace nobody will retry in.

## Decision

A finished task can be dismissed as a stopped or failed one can, with
the same effect: its runner deletes its record, spool and workspace,
and the task takes no more follow-ups. Whether its branch was pushed
does not matter.

The server keeps the task row, its events and its log after dismissal,
so the Tasks page and the task's page still show what it did. The task
page marks it dismissed and offers no actions.

A running or queued task cannot be dismissed; stop it first.

## Consequences

- Dismissing a finished task deletes its workspace. If its branch was
  pushed, the forge has what remains of its work; if not, the work in
  the workspace is lost.
- A retry or a comment on a dismissed finished task is refused with
  "dismissed".
- The server's database still grows with every task; dismissal frees
  runner disks only.
