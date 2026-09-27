-- 0042_create_events.sql
-- Activity events for all tenants. One row per event, subtype-specific data in payload.

BEGIN;

CREATE TYPE event_actor_kind AS ENUM ('user', 'api_key', 'system', 'integration');

CREATE TABLE events (
    id                  uuid            PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           bigint          NOT NULL REFERENCES tenants (id),
    event_type          text            NOT NULL,
    payload_schema_rev  smallint        NOT NULL DEFAULT 1,
    occurred_at         timestamptz     NOT NULL,
    recorded_at         timestamptz     NOT NULL DEFAULT now(),
    source_service      text            NOT NULL,
    actor_kind          event_actor_kind NOT NULL,
    actor_id            text,
    subject_type        text            NOT NULL,
    subject_id          text            NOT NULL,
    correlation_id      uuid,
    causation_id        uuid,
    payload             jsonb           NOT NULL DEFAULT '{}'::jsonb,
    retention_class     text            NOT NULL DEFAULT 'standard',
    redacted_at         timestamptz,
    CONSTRAINT events_event_type_format
        CHECK (event_type ~ '^[a-z_]+\.[a-z_]+$'),
    CONSTRAINT events_retention_class_valid
        CHECK (retention_class IN ('standard', 'extended', 'legal_hold')),
    CONSTRAINT events_payload_is_object
        CHECK (jsonb_typeof(payload) = 'object')
);

CREATE INDEX ix_events_tenant_occurred
    ON events (tenant_id, occurred_at DESC);

CREATE INDEX ix_events_tenant_type_occurred
    ON events (tenant_id, event_type, occurred_at DESC);

CREATE INDEX ix_events_subject
    ON events (tenant_id, subject_type, subject_id, occurred_at DESC);

CREATE INDEX ix_events_correlation
    ON events (correlation_id)
    WHERE correlation_id IS NOT NULL;

CREATE INDEX ix_events_payload_gin
    ON events USING gin (payload jsonb_path_ops);

CREATE INDEX ix_events_retention_sweep
    ON events (retention_class, recorded_at)
    WHERE redacted_at IS NULL;

COMMIT;
