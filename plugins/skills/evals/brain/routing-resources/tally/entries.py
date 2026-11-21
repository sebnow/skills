"""Parsing of time log lines."""

from dataclasses import dataclass
from datetime import date, datetime


@dataclass(frozen=True)
class Entry:
    day: date
    start: datetime
    end: datetime
    client: str
    task: str

    @property
    def minutes(self) -> int:
        return int((self.end - self.start).total_seconds() // 60)


class ParseError(ValueError):
    pass


def parse_line(line: str) -> Entry:
    """Parses `YYYY-MM-DD HH:MM-HH:MM client task text`."""
    parts = line.split(maxsplit=3)
    if len(parts) < 3:
        raise ParseError(f"too few fields: {line!r}")
    day_text, span, client = parts[0], parts[1], parts[2]
    task = parts[3] if len(parts) == 4 else ""
    try:
        day = date.fromisoformat(day_text)
        start_text, end_text = span.split("-")
        start = datetime.combine(day, datetime.strptime(start_text, "%H:%M").time())
        end = datetime.combine(day, datetime.strptime(end_text, "%H:%M").time())
    except ValueError as err:
        raise ParseError(f"bad date or time in {line!r}") from err
    if end <= start:
        raise ParseError(f"end before start: {line!r}")
    return Entry(day, start, end, client, task)


def parse_log(text: str) -> list[Entry]:
    """Parses every non-blank, non-comment line of a log."""
    entries = []
    for line in text.splitlines():
        line = line.strip()
        if not line or line.startswith("#"):
            continue
        entries.append(parse_line(line))
    return entries
