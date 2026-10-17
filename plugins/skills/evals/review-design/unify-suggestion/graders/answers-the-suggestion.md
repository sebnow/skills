---
type: llm
focus: last_message
weight: 1
---
The user asked how to remove the duplication in pantry's three cron jobs effectively. They are `carts.expire_idle` (close carts idle longer than `carts.HOLD` and return their stock), `accounts.purge_lapsed` (delete signups not verified within `accounts.VERIFY_WINDOW`), and `reminders.remind_dormant` (email customers whose last order is older than `reminders.DORMANT_AFTER`, once per spell, recording each reminder before sending it), run from `src/pantry/cron.py`. They share a shape (a time window, an SQL predicate, an action on the matching rows) but no rule: each window is a separate business decision owned by a different module, each predicate filters different columns, and each action differs. The operational plumbing is already shared in `cron.run`. An experienced reviewer answers the point about duplication: unifying the jobs is not worth it, with the reason.

PASS if the report addresses the duplication the user sees, recommends against unifying the jobs' logic, and gives a reason grounded in this code: for example that the windows, predicates, or actions encode separate decisions that change independently, that a shared abstraction would couple them or need options for each difference, or that the shared part is already in `cron.run`.

FAIL if the report does not address the duplication the user sees, agrees to unify the jobs' logic (a shared sweep helper, base class, policy type, or table of policies), leaves the decision open without a recommendation, or rejects it with only a generic reason (YAGNI, premature abstraction, keep it simple) that names nothing in this code.
