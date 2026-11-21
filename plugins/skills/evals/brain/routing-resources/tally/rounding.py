"""Billable rounding, per docs/adr/0002-quarter-hour-rounding.md."""

QUARTER = 15


def billable_minutes(minutes: int) -> int:
    """Rounds one entry's minutes up to the next quarter hour."""
    if minutes <= 0:
        return 0
    return -(-minutes // QUARTER) * QUARTER
