---
type: llm
focus: last_message
weight: 1
---
The package renders overdue-loan reminder emails. Two tests in `tests/test_email.py` only restate the implementation: `test_greeting_comes_before_items_then_total_then_sign_off` asserts the order of sections that the `TEMPLATE` string fixes, and `test_body_matches_template` compares the body against the template's exact text. Neither checks a decision the code makes; each only repeats what the template already says, so it fails on every intended wording or layout change and catches nothing else. Pinning the whole email as a literal expected string, snapshot, or golden file restates the template the same way.

PASS only if the review names at least one of these two tests as having no value because it restates the template or implementation (or equivalently: it only checks what the code dictates, or it creates churn on intended changes), and recommends deleting it or rewriting it to check something logic decides (which loans are listed, the fines, the total), and does not recommend asserting the full email text as a literal, snapshot, or golden file.

FAIL if neither test is criticised for restating the implementation, or if the recommended fix for either test is to pin the full email text (a hand-written literal of the whole body, a snapshot, a golden file), to add more ordering or exact-text assertions, or to improve failure messages.
