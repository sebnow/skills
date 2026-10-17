import secrets
import sqlite3
from datetime import datetime, timedelta

from pantry import db

VERIFY_WINDOW = timedelta(days=7)

# A signup not verified within VERIFY_WINDOW has lapsed: it no longer holds its
# address and can no longer be verified.
LAPSED = "verified_at IS NULL AND created_at < :cutoff"


class EmailTaken(Exception):
    pass


class InvalidToken(Exception):
    pass


def sign_up(conn: sqlite3.Connection, email: str, first_name: str, now: datetime) -> str:
    """Create an unverified account and return its verification token.

    A lapsed signup for the same address is replaced.
    """
    email = email.strip().lower()
    token = secrets.token_urlsafe(32)
    with db.transaction(conn):
        conn.execute(
            f"DELETE FROM accounts WHERE email = :email AND {LAPSED}",
            {"email": email, "cutoff": db.ts(now - VERIFY_WINDOW)},
        )
        try:
            conn.execute(
                "INSERT INTO accounts (email, first_name, verify_token, created_at) VALUES (?, ?, ?, ?)",
                (email, first_name, token, db.ts(now)),
            )
        except sqlite3.IntegrityError:
            raise EmailTaken(email) from None
    return token


def verify(conn: sqlite3.Connection, token: str, now: datetime) -> int:
    """Verify the account the token was issued to and return its id."""
    with db.transaction(conn):
        rows = conn.execute(
            f"""
            UPDATE accounts SET verified_at = :now
            WHERE verify_token = :token AND verified_at IS NULL AND NOT ({LAPSED})
            RETURNING id
            """,
            {"now": db.ts(now), "token": token, "cutoff": db.ts(now - VERIFY_WINDOW)},
        ).fetchall()
    if not rows:
        raise InvalidToken()
    return rows[0][0]


def is_verified(conn: sqlite3.Connection, account_id: int) -> bool:
    row = conn.execute(
        "SELECT verified_at IS NOT NULL FROM accounts WHERE id = ?", (account_id,)
    ).fetchone()
    return bool(row and row[0])


def purge_lapsed(conn: sqlite3.Connection, now: datetime) -> int:
    """Delete lapsed signups and return how many were deleted."""
    with db.transaction(conn):
        return conn.execute(
            f"DELETE FROM accounts WHERE {LAPSED}", {"cutoff": db.ts(now - VERIFY_WINDOW)}
        ).rowcount
