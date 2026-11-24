---
type: tool_used
tool: Write
# Words and marks the style rules ban, in the content of any written file.
# tool_used has no flags, so each word spells out its capital. Inside the JSON
# input a line break is the two characters \n, which leaves no word boundary
# before a word that starts a line; the \\n branch covers that.
input_match: '"content"\s*:\s*"(?:[^"\\]|\\.)*(?:—|\\u2014|(?:\b|\\n)(?:[Pp]lease|[Ss]imply|[Jj]ust|[Ee]asily|[Ee]asy|[Vv]ia|[Uu]tilize)\b|\b[Ee][.]g[.]|\b[Ii][.]e[.]|and/or)'
min: 0
max: 0
---
