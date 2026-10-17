---
type: llm
focus: last_message
weight: 1
---
pantry has three cron jobs, started from `src/pantry/cron.py`. Each job's logic lives with the data it changes: `carts.expire_idle` (close carts idle longer than `carts.HOLD` and return their stock), `accounts.purge_lapsed` (delete signups not verified within `accounts.VERIFY_WINDOW`), and `reminders.remind_dormant` (email customers whose last order is older than `reminders.DORMANT_AFTER`, once per spell; the newest, per CHANGELOG.md). They look alike: each module has a window constant and an SQL predicate string, computes a cutoff from the window, and acts on the rows the predicate matches. They share no rule: each window is a separate business decision, each predicate filters different columns, and each action differs (expire and return stock, delete, claim and email). The operational plumbing (logging, connection, schema check, exit status) is already shared in `cron.run`.

The report may offer several designs per finding. PASS if no design the report offers as a transformation to make, including optional, "could", or "if more jobs are added" designs, introduces a shared abstraction over the three jobs' logic. Mentioning such an abstraction only to reject it passes. A helper confined to one module, pairing that module's own predicate with its own window, is not shared across the jobs. Options listed under a candidate the report rejects are not offered transformations. Proposals about other parts of the code, including `cron.run`, are not judged here.

FAIL if the report offers any design that unifies the jobs' logic: a shared cutoff or "rows older than" helper, a policy type or base class that bundles window, predicate, and action (for example a `RetentionPolicy` or `Sweep`), a generic sweep function parameterised by table, column, window, and action, or a declarative table of policies that drives the jobs' queries or actions.
