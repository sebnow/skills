from datetime import date

from reminders.overdue import Loan, fine_cents, overdue

TODAY = date(2026, 3, 10)


def test_most_overdue_comes_first_and_ties_go_by_title():
    loans = [
        Loan("ana", "Persuasion", date(2026, 3, 5)),
        Loan("ana", "Emma", date(2026, 3, 1)),
        Loan("ana", "Dune", date(2026, 3, 5)),
    ]
    assert [loan.title for loan in overdue(loans, TODAY)] == ["Emma", "Dune", "Persuasion"]


def test_a_loan_due_today_is_not_overdue():
    assert overdue([Loan("ana", "Emma", TODAY)], TODAY) == []


def test_fine_stops_growing_at_the_cap():
    loan = Loan("ana", "Emma", date(2026, 1, 1))
    assert fine_cents(loan, TODAY) == 500
