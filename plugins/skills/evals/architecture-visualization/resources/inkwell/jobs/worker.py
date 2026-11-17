"""Conversion worker.

Each worker coroutine takes one job at a time from the queue, fetches the
source document, runs every requested converter, stores the results, and
records the outcome. Success and final failure both queue a webhook for the
customer, and success records billable pages.
"""

import logging
import random

from inkwell.convert.registry import ConverterRegistry, PermanentError
from inkwell.jobs.models import Job
from inkwell.jobs.queue import JobQueue
from inkwell.billing.meter import UsageMeter
from inkwell.notify.webhooks import WebhookOutbox
from inkwell.storage.blobs import BlobStore
from inkwell.storage.repo import DocumentRepo

log = logging.getLogger(__name__)


def backoff_seconds(attempt: int) -> int:
    """Exponential backoff with jitter: about 10s, 40s, 160s, 640s, ..."""
    base = 10 * 4 ** (attempt - 1)
    return int(base * random.uniform(0.8, 1.2))


class Worker:
    def __init__(
        self,
        queue: JobQueue,
        docs: DocumentRepo,
        blobs: BlobStore,
        converters: ConverterRegistry,
        outbox: WebhookOutbox,
        meter: UsageMeter,
        max_attempts: int,
    ):
        self.queue = queue
        self.docs = docs
        self.blobs = blobs
        self.converters = converters
        self.outbox = outbox
        self.meter = meter
        self.max_attempts = max_attempts

    async def run(self) -> None:
        while True:
            job = await self.queue.next()
            try:
                await self.process(job)
            except PermanentError as exc:
                await self.fail(job, str(exc))
            except Exception as exc:  # transient: renderer down, blob timeout, ...
                log.warning("job %s attempt %s failed: %s", job.id, job.attempts, exc)
                if job.attempts >= self.max_attempts:
                    await self.fail(job, f"gave up after {job.attempts} attempts: {exc}")
                else:
                    await self.queue.retry_later(job, str(exc), backoff_seconds(job.attempts))

    async def process(self, job: Job) -> None:
        doc = await self.docs.get_by_id(job.doc_id)
        if doc is None or doc.status == "deleted":
            async with self.docs.transaction() as tx:
                await self.queue.finish(tx, job)
            return

        await self.docs.mark_processing(doc.id)
        source = await self.blobs.get(doc.source_key)

        result_keys: dict[str, str] = {}
        pages = 0
        for output in doc.outputs:
            converter = self.converters.for_output(output, doc.filename)
            result = await converter.convert(source, doc.filename)
            key = f"{doc.tenant}/{doc.id}/{output}{result.extension}"
            await self.blobs.put(key, result.data, result.content_type)
            result_keys[output] = key
            pages = max(pages, result.pages)

        # Marking the document done, finishing the job and queueing the
        # webhook happen together, so the customer is told exactly once.
        async with self.docs.transaction() as tx:
            await self.docs.mark_done(tx, doc.id, result_keys, pages)
            await self.queue.finish(tx, job)
            if doc.callback_url:
                await self.outbox.add(tx, doc.tenant, doc.callback_url, "document.done", doc.id)
        await self.meter.record(doc.tenant, doc.id, pages)

    async def fail(self, job: Job, error: str) -> None:
        doc = await self.docs.get_by_id(job.doc_id)
        async with self.docs.transaction() as tx:
            await self.docs.mark_failed(tx, job.doc_id, error)
            await self.queue.finish(tx, job)
            if doc and doc.callback_url:
                await self.outbox.add(tx, doc.tenant, doc.callback_url, "document.failed", doc.id)
