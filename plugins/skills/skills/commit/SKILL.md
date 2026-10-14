---
name: commit
description: "Writes commits that follow the repository's own conventions. Use before creating any commit: when asked to commit, 'commit my changes', 'commit this', 'git commit', 'jj commit', 'ready to commit', or when finishing a change that should be committed."
---

# Commit

## A convention guessed from the newest subjects

Write the subject in the convention the history uses, and decide which one
that is from enough subjects to tell conventions apart. The newest few can
all share one prefix that fits more than one convention.

`docs:` is the usual trap. Under Conventional Commits `docs` is a type.
Under Scoped Commits every subject is `<scope>: <description>`, where the
scope is the subsystem or module the commit touches, so `docs:` is just
the scope of a documentation change. A log whose newest entries are all
`docs:` fits both.

Read subjects (`git log -20 --format=%s`, or more) until one fits only one
convention: a type such as `feat`, `fix` or `chore` on a code change, or a
parenthesised scope, points to Conventional Commits; a module,
directory, or `treewide` before the colon points to Scoped Commits. In a
scoped history, the subject names the module the change touches, not what
kind of change it is.
