"""API keys and per-tenant rate limiting, both backed by Postgres."""

import hashlib
import time


class Unauthorized(Exception):
    pass


class RateLimited(Exception):
    def __init__(self, retry_after: int):
        self.retry_after = retry_after


def _digest(key: str) -> str:
    return hashlib.sha256(key.encode()).hexdigest()


class KeyStore:
    def __init__(self, pool):
        self.pool = pool

    def create(self, tenant: str, key: str) -> None:
        with self.pool.connection() as conn:
            conn.execute(
                "INSERT INTO api_keys (digest, tenant) VALUES (%s, %s)",
                (_digest(key), tenant),
            )

    async def tenant_for(self, header: str) -> str:
        if not header.startswith("Bearer ink_"):
            raise Unauthorized()
        key = header.removeprefix("Bearer ")
        async with self.pool.connection() as conn:
            row = await (
                await conn.execute(
                    "SELECT tenant FROM api_keys WHERE digest = %s AND revoked_at IS NULL",
                    (_digest(key),),
                )
            ).fetchone()
        if row is None:
            raise Unauthorized()
        return row[0]


class RateLimiter:
    """Fixed one-minute window per tenant, counted in the rate_windows table."""

    def __init__(self, pool, per_minute: int):
        self.pool = pool
        self.per_minute = per_minute

    async def check(self, tenant: str) -> None:
        window = int(time.time() // 60)
        async with self.pool.connection() as conn:
            row = await (
                await conn.execute(
                    """
                    INSERT INTO rate_windows (tenant, window_start, hits)
                    VALUES (%s, %s, 1)
                    ON CONFLICT (tenant, window_start)
                    DO UPDATE SET hits = rate_windows.hits + 1
                    RETURNING hits
                    """,
                    (tenant, window),
                )
            ).fetchone()
        if row[0] > self.per_minute:
            raise RateLimited(retry_after=60 - int(time.time()) % 60)
