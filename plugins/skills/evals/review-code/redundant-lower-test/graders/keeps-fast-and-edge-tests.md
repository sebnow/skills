---
type: llm
focus: last_message
weight: 1
---
The crate delivers webhook payloads with retries. Two lower-level tests earn their place. `retries_server_errors_with_backoff_until_accepted` in `src/lib.rs` checks instantly most of what the costly, `#[ignore]`d end-to-end test in `tests/sandbox.rs` checks against the partner's live sandbox with real sleeps; a fast test that covers most of a slow one stays beside it. `backoff_stays_at_the_cap_however_many_attempts` in `src/backoff.rs` checks the 60 s cap and the overflow guard at attempts 64 and `u32::MAX`, which no `deliver` test in the crate reaches.

PASS if the review does not recommend deleting either of these two tests or call either redundant, worthless, or a restatement of the implementation. Suggesting extra cases beside them, or not mentioning them, passes.

FAIL if the review recommends deleting `retries_server_errors_with_backoff_until_accepted` because the end-to-end test covers it, or recommends deleting `backoff_stays_at_the_cap_however_many_attempts`, or calls either one redundant.
