---
type: tool_used
tool: Write
input_match: '"file_path"\s*:\s*"[^"]*/documentation/[^"]+[.]md"'
min: 1
---
The repository keeps its pages in `documentation/` (the `docs_dir` in
`mkdocs.yml`). Passes when at least one new page is written there.
