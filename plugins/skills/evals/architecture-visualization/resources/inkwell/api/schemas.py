"""Request parsing and response shapes for the public API."""

from dataclasses import dataclass

ALLOWED_OUTPUTS = {"pdf", "thumbnail", "pdf_a"}


class ValidationError(Exception):
    pass


@dataclass
class ConversionRequest:
    outputs: list[str]
    callback_url: str | None

    @classmethod
    def parse(cls, form) -> "ConversionRequest":
        outputs = [o.strip() for o in form.get("outputs", "pdf").split(",") if o.strip()]
        unknown = set(outputs) - ALLOWED_OUTPUTS
        if unknown:
            raise ValidationError(f"unknown outputs: {', '.join(sorted(unknown))}")
        callback = form.get("callback_url")
        if callback and not callback.startswith("https://"):
            raise ValidationError("callback_url must be https")
        return cls(outputs=outputs, callback_url=callback)


def document_view(doc, urls: dict[str, str]) -> dict:
    view = {
        "id": doc.id,
        "status": doc.status,
        "filename": doc.filename,
        "created_at": doc.created_at.isoformat(),
        "outputs": doc.outputs,
    }
    if doc.status == "done":
        view["results"] = urls
        view["pages"] = doc.pages
    if doc.status == "failed":
        view["error"] = doc.error
    return view
