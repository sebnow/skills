---
status: accepted
date: 2026-09-02
---

# Projects

## Context

Tasks started for the same repository repeat the same runner pool, branch
and environment. People also want to see one repository's tasks together.

## Decision

A project names a runner pool, a branch and an environment. A task may
belong to one project and starts with its values unless the task says
otherwise. A task's project is fixed when it starts.

The settings page lists projects and edits them. The API mirrors it under
`/v1/projects`.

## Consequences

- A task without a project behaves as before.
- Changing a project's settings changes only tasks started afterwards.
