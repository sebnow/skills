"""Weekly hours per client."""

from collections import defaultdict

from .entries import Entry
from .rounding import billable_minutes


def week_key(entry: Entry) -> str:
    year, week, _ = entry.day.isocalendar()
    return f"{year}-W{week:02d}"


def weekly_totals(entries: list[Entry]) -> dict[str, dict[str, int]]:
    """Returns billable minutes keyed by ISO week, then client."""
    totals: dict[str, dict[str, int]] = defaultdict(lambda: defaultdict(int))
    for entry in entries:
        totals[week_key(entry)][entry.client] += billable_minutes(entry.minutes)
    return {week: dict(clients) for week, clients in totals.items()}


def render(totals: dict[str, dict[str, int]]) -> str:
    lines = []
    for week in sorted(totals):
        lines.append(week)
        for client in sorted(totals[week]):
            hours = totals[week][client] / 60
            lines.append(f"  {client:<20} {hours:6.2f} h")
    return "\n".join(lines) + "\n"
