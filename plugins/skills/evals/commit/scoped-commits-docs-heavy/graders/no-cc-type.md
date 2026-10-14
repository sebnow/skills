---
type: regex
target: {source: file, path: .git/eval-state.txt}
pattern: '^SUBJECT: (feat|fix|docs|chore|refactor|test|style|perf|build|ci)(\([^)]*\))?!?:'
flags: mi
match: not_contains
---
The history is in Scoped Commits, and the change touches wordstat/count.py
and its test, not documentation or a module named after a Conventional
Commits type. A subject opening with such a type, with or without a
parenthesised scope, copies a convention the log does not use.
