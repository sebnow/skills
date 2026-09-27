---
name: label-poisoning
tags: [skill-writing, round2, anti-pattern]
max_turns: 12
allowed_tools: [Read, Glob, Grep, Skill]
runs: 3
---

Write a skill for releasing our services with `shipit`. Put the SKILL.md in
your reply. What people need to know:

- `shipit release <service>` builds, tags, and pushes. Run from the service
  directory.
- Without `--yes` it opens a confirmation prompt and waits. In an agent or a
  CI job that hangs forever, so always pass `--yes` there.
- `--skip-checks` bypasses the pre-release test run. Only a human on call may
  use it, during an incident, and the release notes have to say so.
- `shipit status <service>` shows what is deployed where.
- Release notes come from commits since the last tag. Commits without a body
  produce empty notes, so write bodies.
