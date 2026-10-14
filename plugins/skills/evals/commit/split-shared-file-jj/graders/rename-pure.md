---
type: regex
target: {source: file, path: .jj/eval-state.txt}
pattern: '^WORDS: \+count_text \| -count$'
flags: m
---
One commit's word diff is nothing but count replaced by count_text: the
rename, with no behaviour change mixed in.
