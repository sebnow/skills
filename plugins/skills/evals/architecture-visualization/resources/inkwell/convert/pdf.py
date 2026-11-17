"""PDF output, produced by the renderer service.

The renderer is a separate container (LibreOffice behind a small HTTP API,
see deploy/compose.yaml). inkwell never runs LibreOffice in-process: a
malformed document can hang or crash it, and that should not take a worker
down with it.
"""

import httpx

from inkwell.convert import registry

SUPPORTED = (".docx", ".doc", ".odt", ".rtf", ".md", ".txt", ".png", ".jpg")


class RendererClient:
    def __init__(self, base_url: str, timeout: float = 120.0):
        self.http = httpx.AsyncClient(base_url=base_url, timeout=timeout)

    async def to_pdf(self, source: bytes, filename: str, archival: bool) -> tuple[bytes, int]:
        resp = await self.http.post(
            "/convert/pdf",
            files={"file": (filename, source)},
            data={"pdfa": "true" if archival else "false"},
        )
        if resp.status_code == 422:
            raise registry.PermanentError(resp.json().get("reason", "unreadable document"))
        resp.raise_for_status()
        return resp.content, int(resp.headers["x-page-count"])


class PdfConverter:
    def __init__(self, renderer: RendererClient, archival: bool):
        self.renderer = renderer
        self.archival = archival

    async def convert(self, source: bytes, filename: str) -> "registry.ConversionResult":
        if not filename.lower().endswith(SUPPORTED):
            raise registry.PermanentError(f"unsupported file type: {filename}")
        data, pages = await self.renderer.to_pdf(source, filename, self.archival)
        return registry.ConversionResult(data, "application/pdf", ".pdf", pages)
