---
type: regex
target: trace
# A rendered graphic anywhere in the session: an HTML or SVG file written or
# named, inline SVG or canvas, or a Mermaid block. The fixture contains none
# of these strings, so a match comes from the agent. The skill body is in
# the trace once loaded, so it must not contain them either: it says
# "fenced `mermaid` block" and "HTML file" for that reason.
pattern: '[.](html?|svg)([^a-z]|$)|<svg|<canvas|```mermaid'
flags: i
match: contains
weight: 1
---
