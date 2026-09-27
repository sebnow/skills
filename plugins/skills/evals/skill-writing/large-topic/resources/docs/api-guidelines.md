# HTTP API guidelines

These guidelines apply to every HTTP endpoint in our Go services. The
helpers named here live under `internal/`; use them rather than
reimplementing the behaviour.

## Contents

- Resource naming
- Identifiers
- Error format
- Pagination
- Filtering and sorting
- Versioning
- Authentication and authorization
- Idempotency
- Timestamps and money
- Deprecation
- OpenAPI
- Caching and conditional requests
- Bulk and batch operations
- Long-running operations
- Webhooks
- Rate limiting
- Request and response bodies
- File uploads and downloads
- Health, readiness, and metrics
- Logging and tracing
- Testing an endpoint

## Resource naming

Paths are plural nouns in kebab-case: `/purchase-orders`, not
`/purchaseOrder` and not `/purchase_orders`. A path segment names a
collection or a member of one; it never names an operation, a table, or a
service. Nest at most one level: `/customers/{id}/addresses` is fine,
`/customers/{id}/addresses/{id}/verifications` is not. Anything deeper gets
its own top-level resource with a filter: `/address-verifications?address=`.

Actions that are not create, read, update, or delete are POSTs to a verb
sub-resource: `/purchase-orders/{id}/cancel`, `/invoices/{id}/send`. The
verb is imperative and singular. A GET never changes state, including
"mark as read" style side effects; those are POSTs.

Singleton sub-resources use a singular noun: `/customers/{id}/credit-line`.
Do not put a version in the path; see Versioning. Do not put a format in
the path (`.json`); the response is always JSON. Query parameters are
snake_case: `?created_at_gte=`.

Reserved words that must not appear as a path segment because the gateway
routes on them: `admin`, `internal`, `health`, `metrics`, `docs`.

Examples:

    GET  /purchase-orders
    GET  /purchase-orders/{id}
    POST /purchase-orders
    PATCH /purchase-orders/{id}
    POST /purchase-orders/{id}/cancel
    GET  /customers/{id}/addresses

## Identifiers

Public ids are ULIDs prefixed with the resource type and an underscore:
`po_01HZX3K9Q4W5R6T7Y8U9I0P1AS`. The prefix is two to four lowercase
letters. Internal integer ids never leave the service, in any field,
including error details and log lines that may be returned to callers.

The prefix table lives in `internal/ids/prefixes.go`. Add a row there
before using a new prefix; the test in that package fails if a prefix is
reused or if a prefix is used in a handler without a row. Generate ids with
`ids.New(ids.PurchaseOrder)`; never format the string by hand.

Ids are opaque to clients. Do not document their structure, do not accept
them without the prefix, and do not accept an id with the wrong prefix for
the endpoint: `GET /purchase-orders/cu_...` is a 404, not a 400, so that
the error does not confirm which prefixes exist.

External references (a supplier's own order number) go in a field named
`external_id` and are never used as a path parameter.

## Error format

Every non-2xx response is the RFC 9457 problem shape with our extensions:

```json
{
  "type": "https://errors.example.com/purchase-orders/insufficient-credit",
  "title": "Insufficient credit",
  "status": 422,
  "detail": "Order total 1200.00 exceeds credit line 800.00",
  "instance": "/purchase-orders/po_01HZX3K9Q4W5R6T7Y8U9I0P1AS",
  "code": "PO_INSUFFICIENT_CREDIT",
  "trace_id": "4bf92f3577b34da6a3ce929d0e0e4736"
}
```

`code` is stable and documented in `docs/errors.md`; clients switch on it.
Codes are `SCREAMING_SNAKE`, prefixed with the resource abbreviation, and
never reused for a different meaning once published. `type` is a URL that
resolves to that doc entry. `title` is a short constant per code. `detail`
is for humans, may change, and must not contain internal ids, stack traces,
SQL, or hostnames. `trace_id` is the W3C trace id from the request context.

Validation errors are 400 with code `VALIDATION_FAILED` and an `errors`
array:

```json
"errors": [
  {"field": "lines[2].quantity", "code": "MIN", "detail": "must be at least 1"},
  {"field": "currency", "code": "UNSUPPORTED", "detail": "XAU is not supported"}
]
```

Field paths use dotted notation with array indexes in brackets. Return all
validation errors at once, not the first one.

Status mapping:

| Situation | Status |
|---|---|
| Malformed body, unknown field, wrong type | 400 |
| Authentication missing or invalid | 401 |
| Authenticated but not allowed, on a write | 403 |
| Resource does not exist, or exists but the caller may not see it | 404 |
| Method not supported on the path | 405 |
| Conflict with current state (already cancelled) | 409 |
| Semantically invalid (insufficient credit) | 422 |
| Rate limited | 429 with `Retry-After` |
| Unexpected failure | 500, generic detail, full context in logs |
| Dependency unavailable | 503 with `Retry-After` |

A 404 for a resource the caller may not see is 404, not 403, so existence
does not leak. Never return 200 with an error in the body. Never return a
bare string body for an error.

Use `internal/problem` to build responses: `problem.New(ctx, code).Detail(...)`.
It fills `type`, `title`, `status`, and `trace_id` from the code table.

## Pagination

Cursor-based only. Responses carry `next_cursor` (opaque, base64 of the
sort key plus id) and `has_more`:

```json
{
  "data": [...],
  "next_cursor": "eyJzIjoiMjAyNi0wMS0wMVQwMDowMDowMFoiLCJpIjoicG9fMDFI...",
  "has_more": true
}
```

Page size is `limit`, default 50, max 500; a larger value is clamped, not
rejected. `limit=0` is a 400. Cursors expire after 24 hours and an expired
cursor is a 400 with code `CURSOR_EXPIRED`. A cursor from one endpoint
passed to another is a 400 with code `CURSOR_INVALID`. Offset pagination is
not offered; the `offset` parameter is rejected with 400 and code
`OFFSET_UNSUPPORTED`.

The cursor encodes the sort key of the last returned row and its id, so
paging is stable under inserts. The helper is `internal/page`:
`page.Parse(r)` returns the limit and decoded cursor, `page.Encode(row)`
builds the next one. Every list endpoint must have a deterministic total
order; when the sort field is not unique, append id.

Do not include a total count by default. If a client needs one, it is a
separate `GET /purchase-orders/count` endpoint with the same filters, and it
is documented as approximate for collections over 10,000 rows.

## Filtering and sorting

Filters are query parameters named after the field: `?status=open`. Ranges
use `_gte` and `_lte` suffixes: `?created_at_gte=2026-01-01T00:00:00Z`;
`_gt` and `_lt` are not offered. Multiple values are comma-separated and
mean OR: `?status=open,approved`. Filters on different fields are ANDed.
There is no general query language and no `q=` free text unless the
resource has a documented search index.

Filterable fields are declared per resource in a `FilterableFields` set;
an unknown filter parameter is a 400 with code `UNFILTERABLE` and the field
name in `errors`. Do not silently ignore unknown parameters; a typo in a
filter that returns everything is a data leak waiting to happen.

`sort` takes a field name with an optional `-` prefix for descending:
`?sort=-created_at`. One sort field only. Only fields in the resource's
`SortableFields` set are accepted; others are 400 with code `UNSORTABLE`.
The default sort is `-created_at` unless the resource documents another.

Boolean filters accept `true` and `false` only. Enum filters accept the
documented values only. Timestamps in filters follow the timestamp rules.

## Versioning

Version by header: `Api-Version: 2026-03-01`. Missing header means the
oldest supported version, never the newest, so old clients keep working
when a new version ships. An unknown version is a 400 with code
`API_VERSION_UNKNOWN`. Versions are dates, and a new one is minted only for
a breaking change: removing or renaming a field, changing a type, changing
the meaning of a status, tightening validation. Adding fields, endpoints,
enum values, or optional parameters is not breaking and needs no version.

Each version is listed in `docs/versions.md` with the transformation
applied. Each version stays supported for 18 months from the date of the
next one. The middleware in `internal/apiversion` maps the header to a
transform chain that rewrites requests up to the current shape and
responses back down; handlers only ever see the current shape. New
endpoints register their transforms there; an endpoint with no transform
for a version is served as-is.

Never branch on the version inside a handler. If a change cannot be
expressed as a transform, it is a new endpoint.

Responses echo the version served in `Api-Version`, so a client that sent
no header learns what it got.

## Authentication and authorization

Bearer tokens only, in the `Authorization` header. The token is validated
by `internal/auth.Middleware`, which sets the principal on the context.
Handlers never parse the token and never read the header. Query-string
tokens and cookies are not accepted.

Authorization is a call to `auth.Can(ctx, action, resource)` before any
data access. Action names are in `internal/auth/actions.go` and follow
`resource.verb`: `purchase_order.read`, `purchase_order.cancel`. Add the
action there before using it; the test fails on undeclared actions. The
resource argument is the id, or `auth.Collection` for list and create.

Missing authorization is 404 for reads and 403 for writes. List endpoints
filter to what the principal may see rather than failing; a principal with
no access gets an empty page, not an error.

Service-to-service calls use the same bearer scheme with a service
principal; there is no separate internal API and no network-based trust.
A service principal's actions are listed per service in
`deploy/principals/`.

Never log the token, in full or in part. The middleware redacts it from
the request dump on panic.

## Idempotency

POST endpoints that create or charge take an `Idempotency-Key` header. The
key is a client-chosen string up to 255 bytes; a UUID is typical. Keys are
stored for 24 hours with the response status, headers, and body; a repeat
with the same key and body returns the stored response with the header
`Idempotent-Replayed: true`. A repeat with the same key and a different
body is 422 with code `IDEMPOTENCY_KEY_REUSED`. A request without the
header on an endpoint that requires it is 400 with code
`IDEMPOTENCY_KEY_REQUIRED`.

Keys are scoped to the principal; two principals may use the same key.
While the first request is still in flight, a repeat is 409 with code
`IDEMPOTENCY_IN_PROGRESS` and `Retry-After: 1`.

The store is `internal/idempotency`; wrap the handler with
`idempotency.Wrap(handler)`. The wrapper handles storage, replay, and the
in-progress lock. Handlers must not implement their own replay. The body
comparison is on the canonical JSON, so key order does not matter.

PATCH and DELETE are idempotent by construction and do not take the header.
POST action endpoints (`/cancel`) take it when the action has a side effect
outside our database, such as sending an email or calling a payment
provider.

## Timestamps and money

Timestamps are RFC 3339 in UTC with millisecond precision and a `Z`
suffix: `2026-03-01T12:34:56.789Z`. Never an offset, never a bare date
where a time is meant, never a Unix integer. Field names end in `_at`:
`created_at`, `cancelled_at`. Dates without a time (a due date) are
`YYYY-MM-DD` and end in `_on`: `due_on`. Parse input with
`internal/timeparse`, which rejects offsets and truncates to milliseconds.

Money is a string decimal with the currency in a sibling field:

```json
{"amount": "1200.00", "currency": "EUR"}
```

Never a float, never an integer of minor units, never a currency symbol.
Amounts carry exactly the minor-unit digits of the currency: two for EUR,
zero for JPY, three for KWD. Amounts are formatted by
`internal/money.Format` and parsed by `internal/money.Parse`, which know
the minor unit for each currency and reject the wrong precision with code
`AMOUNT_PRECISION`. Negative amounts are allowed only where the field
documents them (credit notes).

Percentages are string decimals with up to four fraction digits:
`"19.0000"`. Quantities are integers unless the unit is fractional, in
which case they are string decimals with the unit in a sibling field.

## Deprecation

A deprecated endpoint or field sends `Deprecation: true` and
`Sunset: <RFC 1123 date>` headers on every response that includes it, and
is listed in `docs/deprecations.md` with the replacement and the sunset
date. Sunset is at least 6 months out. Removal needs a new Api-Version;
until then the deprecated shape stays served to older versions through the
transform chain.

Deprecated fields are also marked `deprecated: true` in the OpenAPI spec
with an `x-sunset` extension. The `internal/deprecation` middleware adds
the headers from the spec; do not set them by hand.

Usage of deprecated endpoints is counted per principal in the
`api_deprecated_calls` metric. A sunset may not pass while any principal
called the endpoint in the previous 30 days without the owning team
contacting them.

## OpenAPI

Every endpoint is described in `api/openapi.yaml`, generated from the
handler annotations by `make openapi`; CI fails if the committed file is
stale. Do not edit the YAML by hand. Annotations live in a comment block
above the handler:

```go
// @route POST /purchase-orders/{id}/cancel
// @summary Cancel a purchase order
// @param id path string true "purchase order id"
// @body CancelRequest
// @response 200 PurchaseOrder
// @error PO_ALREADY_CANCELLED PO_NOT_FOUND
```

Every request and response type has at least one example in the spec, and
the examples are run as tests by `internal/apitest`, so keep them valid;
a change to a field that breaks an example fails CI. Error codes listed
under `@error` must exist in `docs/errors.md`.

Descriptions are full sentences. Field descriptions say the unit, the
format, and whether the field can be null. Enum values are listed with a
one-line meaning each.

## Caching and conditional requests

GET responses for single resources carry an `ETag` derived from the row
version, and honour `If-None-Match` with a 304. Collection responses do not
carry an ETag. PATCH and DELETE accept `If-Match`; a mismatch is 412 with
code `PRECONDITION_FAILED`. Clients that read then write should send it.

`Cache-Control` is `private, no-store` for everything unless the resource
documents otherwise; reference data (currencies, countries) may use
`private, max-age=3600`. Never `public`. The helper is
`internal/httpcache`: `httpcache.ETag(w, row.Version)` and
`httpcache.Match(r, row.Version)`.

## Bulk and batch operations

A bulk endpoint is `POST /purchase-orders/batch` with a `requests` array of
up to 100 items, each an object with the fields of the single-item request.
The response is 207 with a `results` array in the same order, each item
either `{"status": 201, "data": {...}}` or `{"status": 422, "error":
{problem}}`. The batch is not transactional; say so in the description.
A batch with more than 100 items is a 400 with code `BATCH_TOO_LARGE`.

Batch endpoints take an `Idempotency-Key` and the replay is for the whole
batch. Rate limits count each item.

Bulk reads use the list endpoint with an `id` filter:
`GET /purchase-orders?id=po_...,po_...`, up to 100 ids.

## Long-running operations

An operation that takes more than about two seconds returns 202 with an
`Operation` resource:

```json
{
  "id": "op_01HZX...",
  "status": "running",
  "created_at": "...",
  "result": null,
  "error": null
}
```

and a `Location` header pointing at `/operations/{id}`. Clients poll
`GET /operations/{id}`; `status` is `running`, `succeeded`, or `failed`.
On success `result` holds the resource that would have been returned by a
synchronous call; on failure `error` holds a problem object. Operations
are readable for 7 days. The helper is `internal/ops`: `ops.Start(ctx,
func(ctx) (any, error))` returns the operation to write out.

Do not hold the request open for slow work, and do not return 200 with a
partially done result.

## Webhooks

Outbound events go through `internal/events`; a handler calls
`events.Emit(ctx, "purchase_order.cancelled", payload)` after the database
commit, never before. Event names are `resource.past_tense_verb`. The
payload is the current resource shape plus `previous` for updates.

Deliveries carry `Webhook-Id`, `Webhook-Timestamp`, and
`Webhook-Signature` (HMAC-SHA256 over timestamp and body, hex). Receivers
are told to reject timestamps older than five minutes. Retries follow
exponential backoff for 24 hours, then the endpoint is disabled and the
owner emailed. The delivery record is queryable at
`GET /webhook-deliveries?event=`.

Never include a bearer token, an internal id, or a full card number in a
payload.

## Rate limiting

Limits are per principal and per endpoint class: 600 requests per minute
for reads, 120 for writes, 20 for batch, unless the endpoint documents a
different class. Every response carries `RateLimit-Limit`,
`RateLimit-Remaining`, and `RateLimit-Reset` (seconds). Over the limit is
429 with `Retry-After` and code `RATE_LIMITED`. The middleware in
`internal/ratelimit` applies the class from the handler annotation
`@class write`; a handler with no annotation is a 500 in tests.

Limits are enforced at the gateway as well; the in-service limiter is the
one clients see in headers, and the two must agree, so change both in
`deploy/gateway/limits.yaml` and the annotation together.

## Request and response bodies

Request bodies are JSON objects. Unknown fields are a 400 with code
`UNKNOWN_FIELD` and the field name in `errors`; do not ignore them. Null
and absent are different on PATCH: absent means unchanged, null means clear
the field, and a field that cannot be cleared documents that and returns
`VALIDATION_FAILED` on null. On POST, absent optional fields take their
documented default.

Response bodies are JSON objects, never bare arrays; a collection is under
`data`. Field names are snake_case. Empty collections are `[]`, never null.
Optional fields that are unset are present with null, so the shape is
stable. Enum values are lowercase snake_case strings. Booleans are
booleans, never `"true"`.

Nested resources are embedded only one level deep and only when the parent
is meaningless without them (order lines in an order); otherwise the
response carries the id and the client fetches. An `expand` parameter is
not offered.

Request size is capped at 1 MiB by the gateway; file uploads have their own
section.

## File uploads and downloads

Uploads do not go through the JSON API. A client asks for an upload
target with `POST /uploads` (body: `filename`, `content_type`,
`size_bytes`) and receives a presigned URL and an `upload_id`; it PUTs the
bytes to the URL, then references the `upload_id` in the resource that
owns the file. Uploads not referenced within 24 hours are deleted. Allowed
content types and the size cap (25 MiB) live in `internal/uploads/policy.go`.

Downloads are `GET /attachments/{id}/download`, which redirects (302) to
a presigned URL valid for five minutes. The API never streams file bytes
itself.

## Health, readiness, and metrics

Every service exposes `/health` (process is up, always 200) and `/ready`
(dependencies reachable, 503 otherwise) on the internal port only; they
are not authenticated and not routed through the gateway. `/metrics`
serves Prometheus text on the same port. Handlers do not add their own
health checks; register a dependency probe in `internal/ready` instead.

Every endpoint emits `http_request_duration_seconds` with route, method,
and status labels through the standard middleware; do not add per-handler
timers. Business metrics go through `internal/metrics` with names in
`docs/metrics.md`.

## Logging and tracing

Handlers log through the request logger on the context,
`log.From(ctx)`, which already carries `trace_id`, `principal`, `route`,
and `api_version`. Log at `info` for state changes with the resource id,
at `warn` for rejected requests that indicate a client bug, at `error`
only for 5xx. Never log request or response bodies, tokens, or personal
data; use ids.

Every outbound call (database, HTTP, queue) is wrapped by the tracing
helpers in `internal/trace`, so spans nest under the request span. Add a
span for any handler step that can take more than 100 ms, named
`verb.noun`. The `trace_id` in error responses is the same one clients can
quote in support requests, so it must be the request's, not a child span's.

## Testing an endpoint

Each endpoint has three kinds of tests. Handler tests in the package use
`internal/apitest.Call` against an in-memory store and assert on status,
problem code, and body shape. Example tests come from the OpenAPI spec and
run automatically. Contract tests in `tests/contract` run against a
deployed environment and cover authorization, pagination, and versioning
for every endpoint; add the new endpoint to the manifest in
`tests/contract/endpoints.yaml` or CI fails.

Tests assert on `code`, never on `detail`. Tests that need time use the
fake clock from `internal/clock`. Tests that need a principal use
`apitest.As(principal)`; never construct a token.

Before opening a change, run `make api-check`, which regenerates the spec,
runs the example tests, and lints the annotations.
