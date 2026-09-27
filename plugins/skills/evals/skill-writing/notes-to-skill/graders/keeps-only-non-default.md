---
type: llm
focus: last_message
weight: 1
---
The notes mix conventions a competent Go programmer follows unprompted (table-driven tests, descriptive case names, running go test, t.Helper in helpers, keeping tests fast) with three facts specific to this team that an agent cannot know on its own: CI runs with -race and vendored modules so network tests need the integration build tag; internal/clock provides a fake clock to use instead of time.Sleep; and require is used over assert, with the reason.

PASS if the skill text keeps the three team-specific facts and leaves out the generic conventions, or includes a generic convention only with an explicit condition that it stays only if a run without the skill shows the agent getting it wrong.

FAIL if the skill text transcribes the generic conventions as rules without such a condition, or drops any of the three team-specific facts.
