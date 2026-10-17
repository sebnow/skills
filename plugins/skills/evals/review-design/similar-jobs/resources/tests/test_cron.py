import pytest

from pantry import cron, db


def test_job_exits_1_when_the_database_is_missing(tmp_path, monkeypatch):
    path = tmp_path / "pantry.db"
    monkeypatch.setenv("PANTRY_DB", str(path))
    with pytest.raises(SystemExit) as exit:
        cron.purge_unverified()
    assert exit.value.code == 1
    assert not path.exists()


def test_job_exits_1_when_the_schema_is_out_of_date(tmp_path, monkeypatch):
    path = tmp_path / "pantry.db"
    db.connect(str(path), create=True).close()
    monkeypatch.setenv("PANTRY_DB", str(path))
    with pytest.raises(SystemExit) as exit:
        cron.purge_unverified()
    assert exit.value.code == 1


def test_job_runs_on_a_migrated_database(tmp_path, monkeypatch):
    path = tmp_path / "pantry.db"
    db.migrate(db.connect(str(path), create=True))
    monkeypatch.setenv("PANTRY_DB", str(path))
    cron.purge_unverified()
