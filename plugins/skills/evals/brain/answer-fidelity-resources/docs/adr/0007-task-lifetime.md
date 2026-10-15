---
status: accepted
date: 2026-08-21
---

# How long a task lives

## Context

A task runs as a process on a runner. The runner keeps a record of the
task, a spool of log lines not yet acknowledged by the server, and the
workspace the task's checkout lives in. Without a rule for when those go,
runners fill their disks.

## Decision

A task is `queued`, `running`, `paused`, `finished`, `stopped` or
`failed`. Stop ends the process; a stopped task keeps its workspace and
record on the runner so it can be retried.

Retrying a `finished` or `paused` task starts a new process in the same
workspace. Stopping a `finished` or `paused` task ends it without starting
a process, and the runner deletes the session it stored for the task.

## Consequences

- The runner sends log lines as events; the server acknowledges them by
  seq. Once the server has acknowledged every event of a process, the
  runner deletes the spool; the task's record stays until the task is
  dismissed.
- `finished` is not terminal: a retry makes the task `running` again.
- The server never deletes a task row; the Tasks page grows without bound.
