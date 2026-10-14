import pytest

from wordstat.util import human_size, read_text


def test_read_text(tmp_path):
    p = tmp_path / "a.txt"
    p.write_text("héllo", encoding="utf-8")
    assert read_text(str(p)) == "héllo"


@pytest.mark.parametrize(
    "n, want",
    [(0, "0 B"), (1023, "1023 B"), (1536, "1.5 KiB"), (5 * 1024**3, "5.0 GiB")],
)
def test_human_size(n, want):
    assert human_size(n) == want
