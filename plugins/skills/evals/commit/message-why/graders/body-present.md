---
type: regex
target: {source: file, path: .git/eval-state.txt}
pattern: '^BODY: (?!Co-Authored-By)\S'
flags: m
---
The commit has a body beyond trailers.
