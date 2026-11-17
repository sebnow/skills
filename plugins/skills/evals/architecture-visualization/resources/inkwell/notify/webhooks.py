"""Customer webhooks, delivered from an outbox table.

The worker writes a webhook row in the same transaction that finishes the
document; WebhookSender delivers rows separately and retries until the
customer's endpoint answers 2xx or the row is a day old. Each request is
signed with HMAC-SHA256 so customers can verify it came from us.
"""

import asyncio
import hashlib
import hmac
import json
import time

import httpx

GIVE_UP_AFTER_SECONDS = 24 * 3600


class WebhookOutbox:
    def __init__(self, pool):
        self.pool = pool

    async def add(self, tx, tenant: str, url: str, event: str, doc_id: str) -> None:
        payload = json.dumps({"event": event, "document_id": doc_id})
        await tx.execute(
            "INSERT INTO webhook_outbox (tenant, url, payload) VALUES (%s, %s, %s)",
            (tenant, url, payload),
        )

    async def due(self, limit: int = 50):
        async with self.pool.connection() as conn:
            return await (
                await conn.execute(
                    """
                    SELECT id, url, payload, attempts, created_at FROM webhook_outbox
                    WHERE delivered_at IS NULL AND abandoned_at IS NULL AND next_try <= now()
                    ORDER BY next_try LIMIT %s
                    """,
                    (limit,),
                )
            ).fetchall()

    async def delivered(self, row_id: int) -> None:
        async with self.pool.connection() as conn:
            await conn.execute(
                "UPDATE webhook_outbox SET delivered_at = now() WHERE id = %s", (row_id,)
            )

    async def try_again(self, row_id: int, attempts: int, abandon: bool) -> None:
        async with self.pool.connection() as conn:
            if abandon:
                await conn.execute(
                    "UPDATE webhook_outbox SET abandoned_at = now() WHERE id = %s", (row_id,)
                )
            else:
                await conn.execute(
                    "UPDATE webhook_outbox SET attempts = %s, "
                    "next_try = now() + make_interval(secs => %s) WHERE id = %s",
                    (attempts, min(30 * 2 ** attempts, 3600), row_id),
                )

    async def purge_older_than(self, days: int) -> int:
        async with self.pool.connection() as conn:
            cur = await conn.execute(
                "DELETE FROM webhook_outbox WHERE created_at < now() - make_interval(days => %s)",
                (days,),
            )
            return cur.rowcount


class WebhookSender:
    def __init__(self, outbox: WebhookOutbox, secret: str):
        self.outbox = outbox
        self.secret = secret.encode()
        self.http = httpx.AsyncClient(timeout=10.0)

    def sign(self, body: str, ts: int) -> str:
        return hmac.new(self.secret, f"{ts}.{body}".encode(), hashlib.sha256).hexdigest()

    async def run(self) -> None:
        while True:
            for row_id, url, payload, attempts, created_at in await self.outbox.due():
                ts = int(time.time())
                try:
                    resp = await self.http.post(
                        url,
                        content=payload,
                        headers={
                            "content-type": "application/json",
                            "inkwell-timestamp": str(ts),
                            "inkwell-signature": self.sign(payload, ts),
                        },
                    )
                    ok = 200 <= resp.status_code < 300
                except httpx.HTTPError:
                    ok = False
                if ok:
                    await self.outbox.delivered(row_id)
                else:
                    age = ts - created_at.timestamp()
                    await self.outbox.try_again(
                        row_id, attempts + 1, abandon=age > GIVE_UP_AFTER_SECONDS
                    )
            await asyncio.sleep(2)
