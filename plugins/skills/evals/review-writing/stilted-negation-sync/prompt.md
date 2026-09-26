---
name: stilted-negation-sync
tags: [review-writing, stilted]
max_turns: 12
allowed_tools: [Skill, Agent]
runs: 1
---

Review this README section before I merge it, and put the final version in your reply.

```markdown
## How sync works

Each device keeps a full copy of the notes database and syncs with the server
every few minutes. A client has to upload what the server has not seen:
every edit made while offline.

When a note is deleted, the server replaces it with a tombstone. A tombstone
carries no content; it keeps the note's ID and the time of deletion. Clients
that see a tombstone remove the note locally.

Device clocks hold no authority over ordering. The server assigns each change
a sequence number when it arrives, and clients apply changes in that order.
Two devices editing the same note offline leave no trace of which edit came
first, so the server keeps both versions and marks the note as conflicted.

Attachments sync separately and are downloaded on first open.
```
