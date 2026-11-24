---
type: regex
target: trace
# A Write or an Edit call whose file is mkdocs.yml. A Read of the file has a
# different tool name, and the skill body does not hold this call shape.
pattern: '"name":"(?:Write|Edit)","input":\{[^{}]*"file_path":"[^"]*mkdocs[.]yml"'
match: contains
---
The repository's style guide says to add every new page to the `nav` list in
`mkdocs.yml`. Passes when the agent writes or edits that file.
