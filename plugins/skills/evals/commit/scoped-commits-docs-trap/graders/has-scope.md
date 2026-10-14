---
type: regex
target: {source: file, path: .git/eval-state.txt}
pattern: '^SUBJECT: [\w./-]+(, ?[\w./-]+)*( \([^)]*\))?: \S'
flags: m
---
The subject opens with a scope (count:, wordstat:, or several
comma-separated), then a colon and the description, as every commit in the
history does. A plain sentence subject fails.
