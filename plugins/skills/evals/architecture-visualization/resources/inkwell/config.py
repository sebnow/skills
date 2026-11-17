"""Settings read from the environment, once, at process start."""

import os
from dataclasses import dataclass


@dataclass(frozen=True)
class Settings:
    database_url: str
    blob_bucket: str
    blob_endpoint: str | None
    renderer_url: str
    billing_url: str
    billing_token: str
    webhook_signing_secret: str
    worker_concurrency: int
    max_attempts: int
    result_ttl_days: int
    rate_limit_per_minute: int


def load() -> Settings:
    env = os.environ
    return Settings(
        database_url=env["INKWELL_DATABASE_URL"],
        blob_bucket=env.get("INKWELL_BLOB_BUCKET", "inkwell-documents"),
        blob_endpoint=env.get("INKWELL_BLOB_ENDPOINT"),
        renderer_url=env.get("INKWELL_RENDERER_URL", "http://renderer:3000"),
        billing_url=env["INKWELL_BILLING_URL"],
        billing_token=env["INKWELL_BILLING_TOKEN"],
        webhook_signing_secret=env["INKWELL_WEBHOOK_SECRET"],
        worker_concurrency=int(env.get("INKWELL_WORKER_CONCURRENCY", "4")),
        max_attempts=int(env.get("INKWELL_MAX_ATTEMPTS", "5")),
        result_ttl_days=int(env.get("INKWELL_RESULT_TTL_DAYS", "7")),
        rate_limit_per_minute=int(env.get("INKWELL_RATE_LIMIT", "120")),
    )
