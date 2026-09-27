---
name: new-skill-cold
tags: [skill-writing, baseline-first]
max_turns: 12
allowed_tools: [Read, Glob, Grep, Skill]
runs: 3
---

Create a skill for this plugin that teaches the agent to write commit messages
for our jujutsu repos. Two things keep going wrong: the body says what changed
instead of why, and the `Change-Id:` trailer we rely on for cross-referencing
is missing more often than not. Put whatever you produce in your reply.
