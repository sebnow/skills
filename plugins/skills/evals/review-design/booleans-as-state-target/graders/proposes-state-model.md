---
type: llm
focus: last_message
weight: 1
---
pressroom is a small Go service where authors write articles, other users review them, and authors publish and archive them. The article lifecycle (draft, in review, changes requested, approved, published, archived) is stored as five booleans on `Article`: `Submitted`, `Approved`, `Rejected`, `Published`, `Archived`. Each release in CHANGELOG.md (review in 0.2, rejection in 0.3, archiving in 0.4) added a flag and another guard to the service methods. Every method in `service.go` guards its transition with a hand-written combination of flags (for example `a.Submitted || a.Approved || a.Published || a.Archived`) and sets several flags together; `Article.State` derives the state name with a `switch` whose case order matters; and `view` in `http.go` re-derives from the flags whether the reviewer, the review note, and the publish time apply.

The report may offer several alternative designs for a finding. PASS if at least one offered design replaces the booleans with a single explicit lifecycle state that the service methods, `State`, and the HTTP view all read: a state enum or type with defined transitions, a transition table, a state machine, a sum type or per-state structs, or the State pattern. Keeping data that belongs to some states (reviewer, review note, publish time) beside or inside that state is fine.

FAIL if no offered design removes the boolean flags as the source of truth: for example only extracting the guards into helper methods (such as `canEdit()`), adding comments or tests, adding an invariant check that the flags agree, or adding a status field while keeping the booleans it would replace.
