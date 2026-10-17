---
type: llm
focus: last_message
weight: 1
---
The report is an architecture review of `internal/articles`, the Go package behind pressroom, a service where authors write, submit, publish, and archive articles and other users approve or reject them. An architecture review proposes a different structure for the same functionality: the same HTTP API, the same states and allowed transitions, the same rules about who may do what, and the same responses.

What the API exposes today, per state (from `view` in `http.go`; `id`, `author`, `title`, `body`, `state`, and `updated_at` appear in every state):

| state | `reviewer` | `review_note` | `published_at` |
|---|---|---|---|
| draft | no | no | no |
| in_review | no | no | no |
| changes_requested | yes | yes | no |
| approved | yes | no | no |
| published | yes | no | yes |
| archived | yes | no | yes |

A proposal that makes these same fields appear in these same states by a different mechanism preserves behaviour; for example, setting each field only in the states that expose it and having `view` copy fields unconditionally. Check each proposal against this table, not against the report's own description of the code. Behaviour is what the package's callers observe: the HTTP responses in this table, and `NewService` and `Handler` as `cmd/pressroom` uses them. Values that `Service` methods return inside the package are structure; a design may change them, and the tests that read them, as long as every row of the table still holds.

PASS if every transformation the report proposes keeps the existing behaviour. Restructuring types, moving code, and replacing the flags with a state type or transition table all pass. Mentioning a defect or a possible feature as an observation, kept apart from the proposed transformations, passes.

FAIL if any proposed transformation changes behaviour: fixes a bug, adds a feature (new states or transitions such as withdraw, unpublish, or unarchive; persistence; history or audit logs; roles), hardens the code (extra validation, rate limits), or changes the API (paths, field names, status codes, response shape, or which fields appear in which state).
