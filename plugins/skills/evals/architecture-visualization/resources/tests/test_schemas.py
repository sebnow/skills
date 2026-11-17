import pytest

from inkwell.api.schemas import ConversionRequest, ValidationError


def test_default_output_is_pdf():
    assert ConversionRequest.parse({}).outputs == ["pdf"]


def test_rejects_unknown_output():
    with pytest.raises(ValidationError):
        ConversionRequest.parse({"outputs": "pdf,docx"})


def test_callback_must_be_https():
    with pytest.raises(ValidationError):
        ConversionRequest.parse({"callback_url": "http://example.com/hook"})
