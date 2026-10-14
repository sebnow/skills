from datetime import date

from reminders.overdue import Loan, fine_cents, overdue

TEMPLATE = """Hello {member},

These items are overdue:

{items}

Total fines: {total}

Please return them to the front desk.
Millbrook Lending Library"""


def money(cents: int) -> str:
    return f"${cents // 100}.{cents % 100:02d}"


def render(member: str, loans: list[Loan], today: date) -> str:
    late = overdue(loans, today)
    fines = [fine_cents(loan, today) for loan in late]
    items = "\n".join(
        f"- {loan.title} (due {loan.due.isoformat()}, fine {money(fine)})"
        for loan, fine in zip(late, fines)
    )
    return TEMPLATE.format(member=member, items=items, total=money(sum(fines)))
