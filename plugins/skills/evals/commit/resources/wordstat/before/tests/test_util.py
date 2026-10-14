from wordstat.util import read_text


def test_read_text(tmp_path):
    p = tmp_path / "a.txt"
    p.write_text("héllo", encoding="utf-8")
    assert read_text(str(p)) == "héllo"
