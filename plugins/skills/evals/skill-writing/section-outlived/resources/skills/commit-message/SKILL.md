---
name: commit-message
description: "Writes commit messages for jujutsu repositories. Use when describing a change, writing a commit message, or asked to commit."
---

# Commit Message

## Subject line

Write the subject in the imperative mood and keep it to 50 characters or
fewer, so `jj log` shows it on one line: "Add retry to the export client",
not "Added retry logic to the export client for robustness".

## Change-Id trailer

End the body with a `Change-Id:` trailer carrying the jujutsu change id of the
commit being described: the id you were given, or the output of
`jj log -r @ -T change_id`. Review tooling joins commits to review threads on
this trailer; without it the commit is orphaned from its review.
