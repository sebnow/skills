"""Invoice lines for one client."""

from dataclasses import dataclass

from .entries import Entry


@dataclass(frozen=True)
class Line:
    day: str
    task: str
    hours: float
    amount: float


def tenth_hours(minutes: int) -> float:
    """Rounds minutes up to the next six-minute block, in hours."""
    blocks = -(-minutes // 6)
    return blocks / 10


def invoice_lines(entries: list[Entry], client: str, rate: float) -> list[Line]:
    lines = []
    for entry in sorted(entries, key=lambda e: e.start):
        if entry.client != client:
            continue
        hours = tenth_hours(entry.minutes)
        lines.append(Line(entry.day.isoformat(), entry.task, hours, round(hours * rate, 2)))
    return lines


def render(lines: list[Line]) -> str:
    out = [f"{l.day}  {l.task:<40} {l.hours:5.1f} h  {l.amount:9.2f}" for l in lines]
    total = sum(l.amount for l in lines)
    out.append(f"{'total':<54} {total:9.2f}")
    return "\n".join(out) + "\n"
