import json
import sqlite3
from datetime import datetime, timedelta

from pantry import accounts, db

# How long a cart holds its reserved stock without being touched. An open
# cart idle for longer is closed: it cannot be used, and expire_idle returns
# its stock.
HOLD = timedelta(hours=2)
IDLE = "status = 'open' AND updated_at < :cutoff"


class NotVerified(Exception):
    pass


class CartClosed(Exception):
    pass


class OutOfStock(Exception):
    pass


class EmptyCart(Exception):
    pass


def open_cart(conn: sqlite3.Connection, account_id: int, now: datetime) -> int:
    if not accounts.is_verified(conn, account_id):
        raise NotVerified(account_id)
    with db.transaction(conn):
        cur = conn.execute(
            "INSERT INTO carts (account_id, updated_at) VALUES (?, ?)", (account_id, db.ts(now))
        )
    return cur.lastrowid


def _touch(conn: sqlite3.Connection, cart_id: int, now: datetime) -> None:
    """Mark an open cart as used now. Call inside a transaction.

    Raises CartClosed if the cart was ordered, expired, or idle for longer than
    HOLD, whether or not expire_idle() has run for it yet.
    """
    cur = conn.execute(
        f"UPDATE carts SET updated_at = :now WHERE id = :id AND status = 'open' AND NOT ({IDLE})",
        {"now": db.ts(now), "id": cart_id, "cutoff": db.ts(now - HOLD)},
    )
    if cur.rowcount == 0:
        raise CartClosed(cart_id)


def add_item(conn: sqlite3.Connection, cart_id: int, sku: str, quantity: int, now: datetime) -> None:
    """Reserve `quantity` of `sku` for the cart."""
    if quantity < 1:
        raise ValueError(f"quantity must be positive, got {quantity}")
    with db.transaction(conn):
        _touch(conn, cart_id, now)
        cur = conn.execute(
            "UPDATE products SET stock = stock - ? WHERE sku = ? AND stock >= ?",
            (quantity, sku, quantity),
        )
        if cur.rowcount == 0:
            raise OutOfStock(sku)
        conn.execute(
            """
            INSERT INTO cart_items (cart_id, sku, quantity) VALUES (?, ?, ?)
            ON CONFLICT (cart_id, sku) DO UPDATE SET quantity = quantity + excluded.quantity
            """,
            (cart_id, sku, quantity),
        )


def checkout(conn: sqlite3.Connection, cart_id: int, now: datetime) -> int:
    """Turn the cart into an order and return its total in cents."""
    with db.transaction(conn):
        _touch(conn, cart_id, now)
        total, lines = conn.execute(
            """
            SELECT COALESCE(SUM(i.quantity * p.price_cents), 0), COUNT(*)
            FROM cart_items i JOIN products p ON p.sku = i.sku
            WHERE i.cart_id = ?
            """,
            (cart_id,),
        ).fetchone()
        if lines == 0:
            raise EmptyCart(cart_id)
        conn.execute("UPDATE carts SET status = 'ordered' WHERE id = ?", (cart_id,))
        conn.execute(
            "UPDATE accounts SET last_order_at = ? WHERE id = (SELECT account_id FROM carts WHERE id = ?)",
            (db.ts(now), cart_id),
        )
    return total


def expire_idle(conn: sqlite3.Connection, now: datetime) -> int:
    """Mark idle carts expired, return their stock, and return how many there were."""
    with db.transaction(conn):
        idle = conn.execute(
            f"UPDATE carts SET status = 'expired' WHERE {IDLE} RETURNING id",
            {"cutoff": db.ts(now - HOLD)},
        ).fetchall()
        ids = json.dumps([cart_id for (cart_id,) in idle])
        conn.execute(
            """
            UPDATE products SET stock = stock + (
                SELECT SUM(quantity) FROM cart_items
                WHERE sku = products.sku AND cart_id IN (SELECT value FROM json_each(:ids))
            )
            WHERE sku IN (SELECT sku FROM cart_items WHERE cart_id IN (SELECT value FROM json_each(:ids)))
            """,
            {"ids": ids},
        )
    return len(idle)
