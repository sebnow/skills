from dataclasses import dataclass
from datetime import date


@dataclass(frozen=True)
class Loan:
    member: str
    title: str
    due: date


def overdue(loans: list[Loan], today: date) -> list[Loan]:
    """Loans past their due date, most overdue first, ties broken by title."""
    late = [loan for loan in loans if loan.due < today]
    return sorted(late, key=lambda loan: (loan.due, loan.title))


def fine_cents(loan: Loan, today: date, per_day: int = 25, cap: int = 500) -> int:
    days = (today - loan.due).days
    if days <= 0:
        return 0
    return min(days * per_day, cap)
