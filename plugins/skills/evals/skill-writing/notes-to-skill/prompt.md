---
name: notes-to-skill
tags: [skill-writing, non-default-only]
max_turns: 12
allowed_tools: [Read, Glob, Grep, Skill]
runs: 3
---

Turn these notes from our team wiki into a skill for the plugin. Put the
SKILL.md in your reply.

Go testing conventions

- Use table-driven tests.
- Name test cases descriptively.
- Run `go test ./...` before pushing.
- Use `t.Helper()` in test helpers.
- CI runs with `-race` and `GOFLAGS=-mod=vendor`. Tests that touch the network
  fail there unless they sit behind `//go:build integration`.
- Keep tests fast.
- Do not use `time.Sleep` in tests. `internal/clock` has a fake clock; inject it.
- Assert with `testify/require`, not `assert`. A failed require stops the test,
  so a later panic does not hide the real failure.
