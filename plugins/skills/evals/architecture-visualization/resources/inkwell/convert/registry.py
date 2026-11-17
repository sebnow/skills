"""Chooses a converter for each requested output."""

from dataclasses import dataclass
from typing import Protocol

from inkwell.convert.images import ThumbnailConverter
from inkwell.convert.pdf import RendererClient, PdfConverter


class PermanentError(Exception):
    """The document can never be converted; retrying will not help."""


@dataclass
class ConversionResult:
    data: bytes
    content_type: str
    extension: str
    pages: int


class Converter(Protocol):
    async def convert(self, source: bytes, filename: str) -> ConversionResult: ...


class ConverterRegistry:
    def __init__(self, by_output: dict[str, Converter]):
        self.by_output = by_output

    def for_output(self, output: str, filename: str) -> Converter:
        try:
            return self.by_output[output]
        except KeyError:
            raise PermanentError(f"no converter for {output!r} ({filename})")


def build_registry(renderer_url: str) -> ConverterRegistry:
    renderer = RendererClient(renderer_url)
    return ConverterRegistry(
        {
            "pdf": PdfConverter(renderer, archival=False),
            "pdf_a": PdfConverter(renderer, archival=True),
            "thumbnail": ThumbnailConverter(renderer),
        }
    )
