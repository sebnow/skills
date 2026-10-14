---
type: regex
target: {source: file, path: .jj/eval-state.txt}
pattern: '^SUBJECT: '
flags: m
match: "count:3"
---
The working copy holds three independent changes, so the run ends with
exactly three new commits. One bundled commit, or a change split from its
own tests, gives another count.
