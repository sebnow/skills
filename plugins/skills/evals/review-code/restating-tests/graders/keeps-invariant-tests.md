---
type: llm
focus: last_message
weight: 1
---
The package renders overdue-loan reminder emails. Four tests check behaviour that logic decides, not the template: `test_most_overdue_comes_first_and_ties_go_by_title` (an order produced by `sorted` with a key, not by layout), `test_a_loan_due_today_is_not_overdue` (a boundary), `test_fine_stops_growing_at_the_cap` (a cap), and `test_total_is_the_sum_of_capped_fines` (arithmetic across loans). These earn their place.

PASS if the review does not recommend deleting any of these four tests or call them worthless, brittle, or a restatement of the implementation. Suggesting extra cases beside them, or not mentioning them, passes.

FAIL if the review recommends deleting any of the four, or criticises them as restating the implementation — for example rejecting every order assertion including the sort test, or rejecting all tests in a file wholesale.
