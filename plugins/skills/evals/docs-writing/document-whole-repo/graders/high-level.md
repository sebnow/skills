---
type: llm
focus: trace
weight: 1
---
The agent wrote documentation for inkwell, a document conversion service in Python: an HTTP API with API keys and per-tenant rate limits, a Postgres job queue, conversion workers that call an external renderer, object storage, a webhook outbox and sender, a usage meter that reports to a billing provider, and a nightly janitor. Judge only the Markdown pages the agent wrote (the Write and Edit calls), not its chat replies.

The pages should explain the system at the level of its architecture and its use: parts, their responsibilities, boundaries, flows, lifecycles, and the reasons for them, plus what an API consumer or operator does. Details a reader can get easily from the code should stay out.

PASS if the pages are built from such parts and flows, and code-level detail is rare: at most an occasional command, endpoint, or HTTP field that a user of the API needs to type or send.

FAIL if the pages repeat what the code says: source file paths or a module-by-module tour of the package layout, class or function names, SQL or table schemas, environment variable lists with defaults, constants such as retry delays or size limits copied from the code, or a table of every configuration value.
