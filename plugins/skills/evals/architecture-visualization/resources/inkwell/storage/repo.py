"""Document records in Postgres."""

import contextlib
import json
from dataclasses import dataclass
from datetime import datetime


@dataclass
class Document:
    id: str
    tenant: str
    filename: str
    status: str  # queued | processing | done | failed | deleted
    source_key: str
    outputs: list[str]
    callback_url: str | None
    result_keys: dict[str, str]
    pages: int
    error: str | None
    created_at: datetime


_COLUMNS = (
    "id, tenant, filename, status, source_key, outputs, callback_url, "
    "result_keys, pages, error, created_at"
)


class DocumentRepo:
    def __init__(self, pool):
        self.pool = pool

    @contextlib.asynccontextmanager
    async def transaction(self):
        async with self.pool.connection() as conn:
            async with conn.transaction():
                yield conn

    async def insert(self, tx, *, doc_id, tenant, filename, source_key, outputs, callback_url):
        await tx.execute(
            """
            INSERT INTO documents (id, tenant, filename, status, source_key, outputs, callback_url)
            VALUES (%s, %s, %s, 'queued', %s, %s, %s)
            """,
            (doc_id, tenant, filename, source_key, json.dumps(outputs), callback_url),
        )

    async def get(self, tenant: str, doc_id: str) -> Document | None:
        return await self._one(
            f"SELECT {_COLUMNS} FROM documents WHERE tenant = %s AND id = %s AND status <> 'deleted'",
            (tenant, doc_id),
        )

    async def get_by_id(self, doc_id: str) -> Document | None:
        return await self._one(f"SELECT {_COLUMNS} FROM documents WHERE id = %s", (doc_id,))

    async def mark_processing(self, doc_id: str) -> None:
        async with self.pool.connection() as conn:
            await conn.execute(
                "UPDATE documents SET status = 'processing' WHERE id = %s", (doc_id,)
            )

    async def mark_done(self, tx, doc_id: str, result_keys: dict, pages: int) -> None:
        await tx.execute(
            "UPDATE documents SET status = 'done', result_keys = %s, pages = %s, "
            "finished_at = now() WHERE id = %s",
            (json.dumps(result_keys), pages, doc_id),
        )

    async def mark_failed(self, tx, doc_id: str, error: str) -> None:
        await tx.execute(
            "UPDATE documents SET status = 'failed', error = %s, finished_at = now() WHERE id = %s",
            (error, doc_id),
        )

    async def mark_deleted(self, doc_id: str) -> None:
        async with self.pool.connection() as conn:
            await conn.execute("UPDATE documents SET status = 'deleted' WHERE id = %s", (doc_id,))

    async def expired(self, ttl_days: int, limit: int = 500) -> list[Document]:
        async with self.pool.connection() as conn:
            rows = await (
                await conn.execute(
                    f"SELECT {_COLUMNS} FROM documents "
                    "WHERE status IN ('done', 'failed') "
                    "AND finished_at < now() - make_interval(days => %s) LIMIT %s",
                    (ttl_days, limit),
                )
            ).fetchall()
        return [self._row(r) for r in rows]

    async def _one(self, sql, params) -> Document | None:
        async with self.pool.connection() as conn:
            row = await (await conn.execute(sql, params)).fetchone()
        return self._row(row) if row else None

    @staticmethod
    def _row(r) -> Document:
        return Document(
            id=r[0], tenant=r[1], filename=r[2], status=r[3], source_key=r[4],
            outputs=r[5], callback_url=r[6], result_keys=r[7] or {}, pages=r[8] or 0,
            error=r[9], created_at=r[10],
        )
