from datetime import date

from reminders.email import TEMPLATE, render
from reminders.overdue import Loan

TODAY = date(2026, 3, 10)
LOANS = [
    Loan("ana", "Emma", date(2026, 3, 1)),
    Loan("ana", "Dune", date(2026, 3, 8)),
]


def test_greeting_comes_before_items_then_total_then_sign_off():
    body = render("Ana", LOANS, TODAY)
    greeting = body.index("Hello Ana")
    items = body.index("These items are overdue")
    total = body.index("Total fines")
    sign_off = body.index("Millbrook Lending Library")
    assert greeting < items < total < sign_off


def test_body_matches_template():
    body = render("Ana", LOANS, TODAY)
    expected = TEMPLATE.format(
        member="Ana",
        items="- Emma (due 2026-03-01, fine $2.25)\n- Dune (due 2026-03-08, fine $0.50)",
        total="$2.75",
    )
    assert body == expected


def test_total_is_the_sum_of_capped_fines():
    loans = [
        Loan("ana", "Emma", date(2025, 12, 1)),
        Loan("ana", "Dune", date(2026, 3, 8)),
    ]
    assert "Total fines: $5.50" in render("Ana", loans, TODAY)
