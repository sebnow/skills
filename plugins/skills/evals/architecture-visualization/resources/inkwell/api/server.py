"""Public HTTP API.

POST /v1/documents            upload a document and request outputs
GET  /v1/documents/{id}       status and, when done, signed result URLs
DELETE /v1/documents/{id}     delete a document and its results early
"""

import uuid

from starlette.applications import Starlette
from starlette.requests import Request
from starlette.responses import JSONResponse
from starlette.routing import Route

from inkwell.api.auth import KeyStore, RateLimiter, Unauthorized, RateLimited
from inkwell.api.schemas import ConversionRequest, ValidationError, document_view
from inkwell.jobs.queue import JobQueue
from inkwell.storage.blobs import BlobStore
from inkwell.storage.repo import DocumentRepo

MAX_UPLOAD_BYTES = 50 * 1024 * 1024


def build_app(
    keys: KeyStore,
    limiter: RateLimiter,
    docs: DocumentRepo,
    blobs: BlobStore,
    queue: JobQueue,
) -> Starlette:
    async def authenticate(request: Request) -> str:
        tenant = await keys.tenant_for(request.headers.get("authorization", ""))
        await limiter.check(tenant)
        return tenant

    async def create_document(request: Request) -> JSONResponse:
        tenant = await authenticate(request)
        form = await request.form(max_part_size=MAX_UPLOAD_BYTES)
        upload = form["file"]
        conv = ConversionRequest.parse(form)

        doc_id = uuid.uuid4().hex
        source_key = f"{tenant}/{doc_id}/source"
        await blobs.put(source_key, await upload.read(), upload.content_type)

        # The document row and its job are written in one transaction, so a
        # crash between them cannot leave a document no worker will pick up.
        async with docs.transaction() as tx:
            await docs.insert(
                tx,
                doc_id=doc_id,
                tenant=tenant,
                filename=upload.filename,
                source_key=source_key,
                outputs=conv.outputs,
                callback_url=conv.callback_url,
            )
            await queue.enqueue(tx, doc_id=doc_id, kind="convert")
        return JSONResponse({"id": doc_id, "status": "queued"}, status_code=202)

    async def get_document(request: Request) -> JSONResponse:
        tenant = await authenticate(request)
        doc = await docs.get(tenant, request.path_params["doc_id"])
        if doc is None:
            return JSONResponse({"error": "not_found"}, status_code=404)
        urls = {}
        if doc.status == "done":
            urls = {k: blobs.signed_url(v) for k, v in doc.result_keys.items()}
        return JSONResponse(document_view(doc, urls))

    async def delete_document(request: Request) -> JSONResponse:
        tenant = await authenticate(request)
        doc = await docs.get(tenant, request.path_params["doc_id"])
        if doc is None:
            return JSONResponse({"error": "not_found"}, status_code=404)
        await blobs.delete_prefix(f"{tenant}/{doc.id}/")
        await docs.mark_deleted(doc.id)
        return JSONResponse({"id": doc.id, "status": "deleted"})

    async def on_unauthorized(request, exc):
        return JSONResponse({"error": "unauthorized"}, status_code=401)

    async def on_rate_limited(request, exc):
        return JSONResponse(
            {"error": "rate_limited"},
            status_code=429,
            headers={"retry-after": str(exc.retry_after)},
        )

    async def on_invalid(request, exc):
        return JSONResponse({"error": "invalid", "detail": str(exc)}, status_code=422)

    return Starlette(
        routes=[
            Route("/v1/documents", create_document, methods=["POST"]),
            Route("/v1/documents/{doc_id}", get_document, methods=["GET"]),
            Route("/v1/documents/{doc_id}", delete_document, methods=["DELETE"]),
        ],
        exception_handlers={
            Unauthorized: on_unauthorized,
            RateLimited: on_rate_limited,
            ValidationError: on_invalid,
        },
    )
