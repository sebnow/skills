"""Reminders for customers whose last order is older than DORMANT_AFTER.

A customer gets one reminder per dormant spell; ordering again starts a new
spell. Each reminder is recorded before it is sent, so an overlapping run never
sends one twice. A reminder that cannot be written, such as one for a name with
a line break, is skipped without being recorded. A reminder the relay does not
take is released for the next run, unless the relay rejects that message for
good; only a crash between recording and sending can lose one.
"""

import logging
import smtplib
import sqlite3
from datetime import datetime, timedelta

from pantry import db
from pantry.mail import Mailer

DORMANT_AFTER = timedelta(days=60)
DUE = "last_order_at < :cutoff AND (reminded_at IS NULL OR reminded_at < last_order_at)"

log = logging.getLogger(__name__)


def remind_dormant(conn: sqlite3.Connection, mailer: Mailer, now: datetime) -> int:
    """Send the reminders that are due and return how many were sent."""
    cutoff = db.ts(now - DORMANT_AFTER)
    claimed_at = db.ts(now)
    due = conn.execute(
        f"SELECT id, email, first_name FROM accounts WHERE {DUE} ORDER BY id", {"cutoff": cutoff}
    ).fetchall()
    if due:
        # Reach the relay before claiming anyone, so an outage loses no reminders.
        mailer.connect()
    sent = 0
    for account_id, email, first_name in due:
        try:
            msg = mailer.compose("we-miss-you", to=email, first_name=first_name)
        except ValueError as e:
            log.warning("cannot write a reminder for account %d: %s", account_id, e)
            continue
        claimed, previous = _claim(conn, account_id, cutoff, claimed_at)
        if not claimed:
            continue
        try:
            mailer.send(msg)
        except (smtplib.SMTPRecipientsRefused, smtplib.SMTPDataError) as e:
            if _permanent(e):
                log.warning("relay rejected the reminder for account %d: %s", account_id, e)
            else:
                _release(conn, account_id, previous, claimed_at)
                log.warning("relay deferred the reminder for account %d; the next run retries it", account_id)
            continue
        except BaseException:
            _release(conn, account_id, previous, claimed_at)
            raise
        sent += 1
    return sent


def _claim(conn: sqlite3.Connection, account_id: int, cutoff: int, claimed_at: int) -> tuple[bool, int | None]:
    """Record the reminder if it is still due; return whether it was, and the
    previous reminded_at for _release."""
    with db.transaction(conn):
        row = conn.execute(
            f"SELECT reminded_at FROM accounts WHERE id = :id AND {DUE}",
            {"id": account_id, "cutoff": cutoff},
        ).fetchone()
        if row is None:
            return False, None
        conn.execute(
            "UPDATE accounts SET reminded_at = :now WHERE id = :id", {"now": claimed_at, "id": account_id}
        )
    return True, row[0]


def _release(conn: sqlite3.Connection, account_id: int, previous: int | None, claimed_at: int) -> None:
    """Undo a claim whose reminder was not sent, so the next run retries it."""
    with db.transaction(conn):
        conn.execute(
            "UPDATE accounts SET reminded_at = :previous WHERE id = :id AND reminded_at = :claimed",
            {"previous": previous, "id": account_id, "claimed": claimed_at},
        )


def _permanent(e: smtplib.SMTPRecipientsRefused | smtplib.SMTPDataError) -> bool:
    """Whether the relay rejected this one message for good (a 5xx reply)."""
    if isinstance(e, smtplib.SMTPRecipientsRefused):
        return all(code >= 500 for code, _ in e.recipients.values())
    return e.smtp_code >= 500
