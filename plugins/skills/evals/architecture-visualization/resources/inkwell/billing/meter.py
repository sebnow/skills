"""Usage metering.

Billable pages are recorded locally per document and pushed to the billing
provider in batches once a minute. The provider deduplicates on
idempotency_key, so a batch that is sent twice is not billed twice.
"""

import asyncio

import httpx

FLUSH_SECONDS = 60


class UsageMeter:
    def __init__(self, pool, billing_url: str, token: str):
        self.pool = pool
        self.http = httpx.AsyncClient(
            base_url=billing_url, headers={"authorization": f"Bearer {token}"}, timeout=30.0
        )

    async def record(self, tenant: str, doc_id: str, pages: int) -> None:
        async with self.pool.connection() as conn:
            await conn.execute(
                "INSERT INTO usage_events (tenant, doc_id, pages) VALUES (%s, %s, %s) "
                "ON CONFLICT (doc_id) DO NOTHING",
                (tenant, doc_id, pages),
            )

    async def run_flusher(self) -> None:
        while True:
            await asyncio.sleep(FLUSH_SECONDS)
            try:
                await self.flush()
            except httpx.HTTPError:
                pass  # left unsent; the next flush picks the events up again

    async def flush(self) -> None:
        async with self.pool.connection() as conn:
            rows = await (
                await conn.execute(
                    "SELECT id, tenant, doc_id, pages FROM usage_events "
                    "WHERE sent_at IS NULL ORDER BY id LIMIT 1000"
                )
            ).fetchall()
            if not rows:
                return
            events = [
                {"idempotency_key": f"doc-{doc_id}", "customer": tenant, "quantity": pages,
                 "metric": "pages_converted"}
                for _, tenant, doc_id, pages in rows
            ]
            resp = await self.http.post("/v1/meter-events/batch", json={"events": events})
            resp.raise_for_status()
            await conn.execute(
                "UPDATE usage_events SET sent_at = now() WHERE id = ANY(%s)",
                ([r[0] for r in rows],),
            )
