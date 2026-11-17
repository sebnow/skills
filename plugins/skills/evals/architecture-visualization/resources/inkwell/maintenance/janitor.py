"""Nightly cleanup, run from cron.

Results are kept for result_ttl_days after a document finishes. The janitor
deletes their blobs, marks the documents deleted, and purges delivered or
abandoned webhook rows older than thirty days.
"""

import logging

log = logging.getLogger(__name__)


class Janitor:
    def __init__(self, docs, blobs, outbox, ttl_days: int):
        self.docs = docs
        self.blobs = blobs
        self.outbox = outbox
        self.ttl_days = ttl_days

    async def run_once(self) -> None:
        removed = 0
        while batch := await self.docs.expired(self.ttl_days):
            for doc in batch:
                await self.blobs.delete_prefix(f"{doc.tenant}/{doc.id}/")
                await self.docs.mark_deleted(doc.id)
                removed += 1
        purged = await self.outbox.purge_older_than(30)
        log.info("janitor: removed %d documents, purged %d webhook rows", removed, purged)
