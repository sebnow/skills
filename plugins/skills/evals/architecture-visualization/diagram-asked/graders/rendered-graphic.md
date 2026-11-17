---
type: regex
target: trace
# A rendered graphic anywhere in the session: an HTML or SVG file written or
# named, inline SVG or canvas, or a Mermaid block. The fixture contains none
# of these strings, so a match comes from the agent.
pattern: '[.](html?|svg)([^a-z]|$)|<svg|<canvas|```mermaid'
flags: i
match: contains
weight: 1
---
