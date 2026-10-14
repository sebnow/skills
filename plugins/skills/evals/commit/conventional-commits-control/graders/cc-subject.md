---
type: regex
target: {source: file, path: .git/eval-state.txt}
pattern: '^SUBJECT: (feat|fix|docs|chore|refactor|test|style|perf|build|ci)(\([^)]+\))?!?: \S'
flags: m
---
The history is in Conventional Commits, so the subject opens with a type,
an optional parenthesised scope, and a colon, as in fix(count): ...
