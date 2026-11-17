"""A job queue in Postgres.

Workers claim jobs with SELECT ... FOR UPDATE SKIP LOCKED, so any number of
worker processes can share the table without handing one job to two of them.
A claimed job is leased for LEASE_SECONDS; if the worker dies, the lease runs
out and another worker picks the job up again.
"""

import asyncio

from inkwell.jobs.models import Job

LEASE_SECONDS = 300
POLL_SECONDS = 1.0


class JobQueue:
    def __init__(self, pool):
        self.pool = pool

    async def enqueue(self, tx, doc_id: str, kind: str, delay_seconds: int = 0) -> None:
        await tx.execute(
            """
            INSERT INTO jobs (doc_id, kind, run_after)
            VALUES (%s, %s, now() + make_interval(secs => %s))
            """,
            (doc_id, kind, delay_seconds),
        )

    async def claim(self) -> Job | None:
        async with self.pool.connection() as conn:
            row = await (
                await conn.execute(
                    """
                    UPDATE jobs SET leased_until = now() + make_interval(secs => %s),
                                    attempts = attempts + 1
                    WHERE id = (
                        SELECT id FROM jobs
                        WHERE run_after <= now()
                          AND (leased_until IS NULL OR leased_until < now())
                          AND finished_at IS NULL
                        ORDER BY run_after
                        FOR UPDATE SKIP LOCKED
                        LIMIT 1
                    )
                    RETURNING id, doc_id, kind, attempts, run_after, last_error
                    """,
                    (LEASE_SECONDS,),
                )
            ).fetchone()
        return Job(*row) if row else None

    async def next(self) -> Job:
        while True:
            job = await self.claim()
            if job is not None:
                return job
            await asyncio.sleep(POLL_SECONDS)

    async def finish(self, tx, job: Job) -> None:
        await tx.execute("UPDATE jobs SET finished_at = now() WHERE id = %s", (job.id,))

    async def retry_later(self, job: Job, error: str, delay_seconds: int) -> None:
        async with self.pool.connection() as conn:
            await conn.execute(
                """
                UPDATE jobs SET leased_until = NULL, last_error = %s,
                                run_after = now() + make_interval(secs => %s)
                WHERE id = %s
                """,
                (error, delay_seconds, job.id),
            )
