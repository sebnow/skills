---
type: llm
focus: last_message
weight: 1
---
deskbook is a Go service for booking meeting rooms. The three commits under review add two room rules and a dry-run endpoint: a room's `max_minutes` cap (commit 1), `members_only` rooms (commit 2), and `POST /bookings/check` (commit 3). Both rules are decided in the HTTP layer: `book` in `internal/api/api.go` checks them inline, and commit 3 copies the same two checks into `check`. The domain already owns the other rule for a valid booking, the overlap check in `bookings.Store`, and the rooms and their settings live in `internal/rooms`. So the decision "does this room accept this booking from this user" sits in the wrong place: it belongs with the domain (`rooms` or `bookings`), where both endpoints and any later caller would get it from one implementation.

What callers observe today, for both `POST /bookings` and `POST /bookings/check`, checked in this order (from `readBooking`, `book`, and `check` in `internal/api/api.go`):

| order | condition | status | error message |
|---|---|---|---|
| 1 | no `X-Forwarded-User` | 401 | `missing X-Forwarded-User` |
| 2 | body not one JSON object, or unknown field | 400 | `invalid JSON body: …` |
| 2 | body over 64 KiB | 413 | `request body too large` |
| 3 | end not after start | 400 | `end must be after start` |
| 4 | start before now | 400 | `start is in the past` |
| 5 | room not configured | 404 | `unknown room` |
| 6 | room is `members_only` and the user is not in `members` | 403 | `<room name> is for members only` |
| 7 | `max_minutes` > 0 and the booking is longer | 422 | `<room name> can be booked for at most <n> minutes` |
| 8 | overlaps another booking of the room | 409 | `the room is already booked for part of that time` |
| – | accepted | 201 with the booking (`POST /bookings`); 200 `{"ok": true}` and nothing booked (`check`) | |

The config format (`max_minutes`, `members_only`, `members`) is also observable to operators. Check each design against this table, not against the report's own description of the code.

The report may offer several designs. PASS if at least one offered design moves the members-only and length-cap rules out of the `api` package into the domain (for example a method on `rooms.Room` or `rooms.Config`, a function or policy type in `bookings`, or `bookings.Store` applying them in `Add` and in a check method), so both endpoints use one implementation, and that design keeps every row of the table: the same statuses, the same messages, and the same order of checks. Mapping domain errors to these statuses in the HTTP layer is fine.

FAIL if no offered design meets both conditions: for example every design keeps the rules in the `api` package (including extracting one shared helper inside `api`), or every design that moves them changes something in the table or the config format, such as checking overlap before the room rules, changing a status or message, or folding `check` into `POST /bookings` behind a flag.
