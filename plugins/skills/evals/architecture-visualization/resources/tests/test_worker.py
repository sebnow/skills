from inkwell.jobs.worker import backoff_seconds


def test_backoff_grows_by_four():
    assert 8 <= backoff_seconds(1) <= 12
    assert 32 <= backoff_seconds(2) <= 48
    assert 128 <= backoff_seconds(3) <= 192
