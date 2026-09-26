---
name: dismissed-alternative
tags: [review-writing, dismissed]
max_turns: 12
allowed_tools: [Skill, Agent]
runs: 1
---

Review this section of our export design note before I commit it, and put the final version in your reply.

```markdown
## Entry format

A Beancount entry has one single-character flag. `*` ("checked") and `!`
("to be looked at later") are the conventional values. The format accepts a
few other characters, which the project does not use. The exporter writes `*`
for rows the user has reconciled and `!` for the rest.

Lots are booked FIFO. Beancount supports several other booking methods, none
of which the exporter implements. Each sell reduces the oldest open lot first.

Dates are written as `YYYY-MM-DD` in the custodian's local date. Timestamps
are dropped; the database stores them, but the export has no use for them.
```
