---
type: regex
target: files
pattern: '\.go$'
flags: m
match: not_contains
---
The brief asks for no tests and no new files. Fails when the junior created
any Go file, such as a new test for the error path.
