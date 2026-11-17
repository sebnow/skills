"""Entry points. Each process builds only the parts it needs."""

import asyncio
import secrets
import sys

import uvicorn

from inkwell import config
from inkwell.api.auth import KeyStore, RateLimiter
from inkwell.api.server import build_app
from inkwell.billing.meter import UsageMeter
from inkwell.convert.registry import build_registry
from inkwell.jobs.queue import JobQueue
from inkwell.jobs.worker import Worker
from inkwell.maintenance.janitor import Janitor
from inkwell.notify.webhooks import WebhookOutbox, WebhookSender
from inkwell.storage.blobs import BlobStore
from inkwell.storage.db import connect
from inkwell.storage.repo import DocumentRepo


def serve(settings: config.Settings) -> None:
    pool = connect(settings.database_url)
    app = build_app(
        keys=KeyStore(pool),
        limiter=RateLimiter(pool, settings.rate_limit_per_minute),
        docs=DocumentRepo(pool),
        blobs=BlobStore(settings.blob_bucket, settings.blob_endpoint),
        queue=JobQueue(pool),
    )
    uvicorn.run(app, host="0.0.0.0", port=8080)


def worker(settings: config.Settings) -> None:
    pool = connect(settings.database_url)
    blobs = BlobStore(settings.blob_bucket, settings.blob_endpoint)
    outbox = WebhookOutbox(pool)
    w = Worker(
        queue=JobQueue(pool),
        docs=DocumentRepo(pool),
        blobs=blobs,
        converters=build_registry(settings.renderer_url),
        outbox=outbox,
        meter=UsageMeter(pool, settings.billing_url, settings.billing_token),
        max_attempts=settings.max_attempts,
    )
    sender = WebhookSender(outbox, settings.webhook_signing_secret)

    async def main() -> None:
        await asyncio.gather(
            *(w.run() for _ in range(settings.worker_concurrency)),
            sender.run(),
            w.meter.run_flusher(),
        )

    asyncio.run(main())


def janitor(settings: config.Settings) -> None:
    pool = connect(settings.database_url)
    j = Janitor(
        DocumentRepo(pool),
        BlobStore(settings.blob_bucket, settings.blob_endpoint),
        WebhookOutbox(pool),
        ttl_days=settings.result_ttl_days,
    )
    asyncio.run(j.run_once())


def create_key(settings: config.Settings, tenant: str) -> None:
    pool = connect(settings.database_url)
    key = "ink_" + secrets.token_urlsafe(24)
    KeyStore(pool).create(tenant, key)
    print(key)


if __name__ == "__main__":
    s = config.load()
    cmd = sys.argv[1] if len(sys.argv) > 1 else "serve"
    if cmd == "serve":
        serve(s)
    elif cmd == "worker":
        worker(s)
    elif cmd == "janitor":
        janitor(s)
    elif cmd == "create-key":
        create_key(s, sys.argv[sys.argv.index("--tenant") + 1])
    else:
        sys.exit(f"unknown command {cmd!r}")
