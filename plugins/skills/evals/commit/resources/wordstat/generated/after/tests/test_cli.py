import json

from wordstat.cli import main


def test_prints_counts_per_file(tmp_path, capsys):
    p = tmp_path / "a.txt"
    p.write_text("one two three\n", encoding="utf-8")
    assert main([str(p)]) == 0
    assert capsys.readouterr().out == f"      1       3      14 {p}\n"


def test_json_format_prints_one_object_per_file(tmp_path, capsys):
    p = tmp_path / "a.txt"
    p.write_text("one two three\n", encoding="utf-8")
    assert main(["--format", "json", str(p)]) == 0
    assert json.loads(capsys.readouterr().out) == {
        "path": str(p),
        "lines": 1,
        "words": 3,
        "chars": 14,
    }
