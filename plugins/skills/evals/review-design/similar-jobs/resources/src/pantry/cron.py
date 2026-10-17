"""Entry points for the cron jobs. Each job's logic lives with the data it changes."""

import logging
import sqlite3
import sys
from collections.abc import Callable
from contextlib import closing
from datetime import UTC, datetime

from pantry import accounts, carts, db, reminders
from pantry.mail import Mailer

log = logging.getLogger(__name__)


def run(name: str, job: Callable[[sqlite3.Connection, datetime], int]) -> None:
    """Run one job: log to stderr, check the schema, and exit 1 if it fails."""
    logging.basicConfig(level=logging.INFO, format="%(asctime)s %(name)s %(levelname)s %(message)s")
    try:
        with closing(db.connect()) as conn:
            db.check_schema(conn)
            count = job(conn, datetime.now(UTC))
    except Exception:
        log.exception("%s failed", name)
        sys.exit(1)
    log.info("%s done: %d", name, count)


def expire_carts() -> None:
    run("expire-carts", carts.expire_idle)


def purge_unverified() -> None:
    run("purge-unverified", accounts.purge_lapsed)


def remind_dormant() -> None:
    def job(conn: sqlite3.Connection, now: datetime) -> int:
        with Mailer.from_env() as mailer:
            return reminders.remind_dormant(conn, mailer, now)

    run("remind-dormant", job)
