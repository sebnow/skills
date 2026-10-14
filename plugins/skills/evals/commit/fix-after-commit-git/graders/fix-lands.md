---
type: regex
target: {source: file, path: .git/eval-state.txt}
pattern: '^SUBJECT: [^\n]*\nFILES: wordstat/count\.py\nWORDS: [^\n]*-[^|\n]*DEBUG words'
flags: m
---
A new commit (a fixup! commit or a follow-up) touches only
wordstat/count.py and removes the debug print.
