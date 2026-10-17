---
type: llm
focus: last_message
weight: 1
---
deskbook is a Go service for booking meeting rooms. The three commits under review add a `max_minutes` cap and `members_only` rooms, decided inline in the `book` HTTP handler in `internal/api/api.go`, and a dry-run `POST /bookings/check` that copies the same two checks into `check`. The rules belong with the domain: they decide which bookings a room accepts, and `bookings.Store` already enforces the other such rule, overlap, so the rules for a valid booking are split between the store and two HTTP handlers. Commit 3 shows the cost: it had to copy the checks, the copies can drift apart, and any further way in (a CLI, an import, another endpoint) would have to copy them again or skip them.

PASS if the report proposes moving these rules into the domain (`rooms`, `bookings`, or the store) and gives a reason grounded in this code. A domain design offered as an alternative counts even if the report recommends another. Reasons that count: the rules are about rooms and bookings rather than HTTP; the store already owns the overlap rule, so the booking rules are split across layers; or the copy in commit 3, with the risk that the copies drift or that another entry point bypasses them.

FAIL if the report does not propose moving the rules out of the HTTP layer, or proposes it with only a generic reason (separation of concerns, layering, clean architecture, single responsibility) that names nothing in this code.
