---
type: regex
target: last_message
pattern: '^(?:(?!\n(?:#{1,6}[ \t]+|\*\*)[^\n]{0,30}(?:evidence|context|details)[^\n]{0,30}\n)[\s\S])*?\n[ \t]*(?:[a-z][.)]|\([a-z]\))[ \t]'
flags: i
match: contains
---
Passes when, before the evidence heading, at least one decision offers its
options as a lettered list ("a." or "(a)" starting a line), so long options
read as their own block rather than inline in the question. The buried-ask
decisions have options a clause long, so a reply that lists none of them
fails. Does not check that every such decision lists
its options.
