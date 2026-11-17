---
type: regex
target: trace
# Box-drawing characters, or boxes drawn with a plus and three or more dashes:
# a text diagram in the reply or in a file. The fixture contains none.
pattern: '[┌┐└┘├┤┬┴┼│─║╔╗╚╝═]|[+]-{3,}'
match: not_contains
weight: 1
---
