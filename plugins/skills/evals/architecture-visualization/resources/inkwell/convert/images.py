"""Thumbnail of the first page: rendered to PDF by the renderer, then
rasterised and scaled here."""

import io

from PIL import Image

from inkwell.convert import registry

THUMB_WIDTH = 320


class ThumbnailConverter:
    def __init__(self, renderer):
        self.renderer = renderer

    async def convert(self, source: bytes, filename: str) -> "registry.ConversionResult":
        pdf, pages = await self.renderer.to_pdf(source, filename, archival=False)
        first_page = await self._rasterise_first_page(pdf)
        img = Image.open(io.BytesIO(first_page))
        ratio = THUMB_WIDTH / img.width
        img = img.resize((THUMB_WIDTH, int(img.height * ratio)))
        out = io.BytesIO()
        img.save(out, format="PNG")
        return registry.ConversionResult(out.getvalue(), "image/png", ".png", pages)

    async def _rasterise_first_page(self, pdf: bytes) -> bytes:
        resp = await self.renderer.http.post(
            "/convert/png", files={"file": ("doc.pdf", pdf)}, data={"page": "1"}
        )
        resp.raise_for_status()
        return resp.content
