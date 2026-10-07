---
type: regex
target: last_message
pattern: 'check\.sh: all 6 checks passed'
match: contains
---
check.sh's last line of output is "check.sh: all 6 checks passed". Fails when
the report summarises the result ("all checks pass") instead of pasting the
output.
