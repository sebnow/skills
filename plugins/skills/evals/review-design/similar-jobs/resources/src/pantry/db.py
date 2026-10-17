import os
import sqlite3
from contextlib import contextmanager
from datetime import datetime
from pathlib import Path

# One entry per schema version; never edit an entry that has shipped.
MIGRATIONS = [
    """
    CREATE TABLE accounts (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        email TEXT NOT NULL UNIQUE,
        first_name TEXT NOT NULL,
        verify_token TEXT NOT NULL UNIQUE,
        created_at INTEGER NOT NULL,
        verified_at INTEGER
    );
    CREATE TABLE products (
        sku TEXT PRIMARY KEY,
        name TEXT NOT NULL,
        price_cents INTEGER NOT NULL CHECK (price_cents >= 0),
        stock INTEGER NOT NULL CHECK (stock >= 0)
    );
    CREATE TABLE carts (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        account_id INTEGER NOT NULL REFERENCES accounts(id),
        status TEXT NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'ordered', 'expired')),
        updated_at INTEGER NOT NULL
    );
    CREATE TABLE cart_items (
        cart_id INTEGER NOT NULL REFERENCES carts(id),
        sku TEXT NOT NULL REFERENCES products(sku),
        quantity INTEGER NOT NULL CHECK (quantity > 0),
        PRIMARY KEY (cart_id, sku)
    );
    """,
    "ALTER TABLE accounts ADD COLUMN last_order_at INTEGER;",
    "CREATE INDEX carts_open_by_age ON carts (updated_at) WHERE status = 'open';",
    "CREATE INDEX accounts_unverified_by_age ON accounts (created_at) WHERE verified_at IS NULL;",
    """
    ALTER TABLE accounts ADD COLUMN reminded_at INTEGER;
    CREATE INDEX accounts_by_last_order ON accounts (last_order_at);
    """,
]


class SchemaOutOfDate(Exception):
    pass


def connect(path: str | None = None, *, create: bool = False) -> sqlite3.Connection:
    """Open the database. Unless create is set, a missing file is an error, not a
    new empty database."""
    path = path or os.environ.get("PANTRY_DB", "pantry.db")
    if path == ":memory:":
        conn = sqlite3.connect(path, isolation_level=None, timeout=10)
    else:
        mode = "rwc" if create else "rw"
        uri = f"{Path(path).absolute().as_uri()}?mode={mode}"
        conn = sqlite3.connect(uri, uri=True, isolation_level=None, timeout=10)
    conn.execute("PRAGMA foreign_keys = ON")
    return conn


def migrate(conn: sqlite3.Connection) -> None:
    (current,) = conn.execute("PRAGMA user_version").fetchone()
    for version, script in enumerate(MIGRATIONS[current:], start=current + 1):
        try:
            conn.executescript(f"BEGIN IMMEDIATE; {script} PRAGMA user_version = {version}; COMMIT;")
        except sqlite3.Error:
            if conn.in_transaction:
                conn.execute("ROLLBACK")
            raise


def check_schema(conn: sqlite3.Connection) -> None:
    (current,) = conn.execute("PRAGMA user_version").fetchone()
    if current != len(MIGRATIONS):
        raise SchemaOutOfDate(f"schema version {current}, expected {len(MIGRATIONS)}; run pantry-migrate")


@contextmanager
def transaction(conn: sqlite3.Connection):
    conn.execute("BEGIN IMMEDIATE")
    try:
        yield
    except BaseException:
        conn.execute("ROLLBACK")
        raise
    conn.execute("COMMIT")


def ts(dt: datetime) -> int:
    """Seconds since the epoch, the form every time is stored and compared in."""
    if dt.tzinfo is None:
        raise ValueError(f"naive datetime {dt!r}")
    return int(dt.timestamp())


def main() -> None:
    migrate(connect(create=True))
