---
name: change-id-trailer
tags: [commit-message]
runs: 3
max_turns: 8
allowed_tools: [Read, Glob, Grep, Skill]
---

Write the commit message for the change in your working directory. The diff
is in CHANGES.diff and the jujutsu change id is in CHANGE_ID. Put the message
in your reply.
