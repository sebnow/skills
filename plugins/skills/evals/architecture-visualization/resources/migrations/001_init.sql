CREATE TABLE api_keys (
    digest      text PRIMARY KEY,
    tenant      text NOT NULL,
    created_at  timestamptz NOT NULL DEFAULT now(),
    revoked_at  timestamptz
);

CREATE TABLE rate_windows (
    tenant        text NOT NULL,
    window_start  bigint NOT NULL,
    hits          int NOT NULL,
    PRIMARY KEY (tenant, window_start)
);

CREATE TABLE documents (
    id            text PRIMARY KEY,
    tenant        text NOT NULL,
    filename      text NOT NULL,
    status        text NOT NULL,
    source_key    text NOT NULL,
    outputs       jsonb NOT NULL,
    callback_url  text,
    result_keys   jsonb,
    pages         int,
    error         text,
    created_at    timestamptz NOT NULL DEFAULT now(),
    finished_at   timestamptz
);
CREATE INDEX documents_expiry ON documents (finished_at) WHERE status IN ('done', 'failed');

CREATE TABLE jobs (
    id            bigserial PRIMARY KEY,
    doc_id        text NOT NULL REFERENCES documents (id),
    kind          text NOT NULL,
    attempts      int NOT NULL DEFAULT 0,
    run_after     timestamptz NOT NULL DEFAULT now(),
    leased_until  timestamptz,
    finished_at   timestamptz,
    last_error    text
);
CREATE INDEX jobs_ready ON jobs (run_after) WHERE finished_at IS NULL;

CREATE TABLE webhook_outbox (
    id            bigserial PRIMARY KEY,
    tenant        text NOT NULL,
    url           text NOT NULL,
    payload       text NOT NULL,
    attempts      int NOT NULL DEFAULT 0,
    next_try      timestamptz NOT NULL DEFAULT now(),
    created_at    timestamptz NOT NULL DEFAULT now(),
    delivered_at  timestamptz,
    abandoned_at  timestamptz
);

CREATE TABLE usage_events (
    id         bigserial PRIMARY KEY,
    tenant     text NOT NULL,
    doc_id     text NOT NULL UNIQUE,
    pages      int NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    sent_at    timestamptz
);
