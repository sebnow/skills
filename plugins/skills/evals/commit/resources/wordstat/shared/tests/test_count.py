from wordstat.count import Counts, count_text


def test_counts_a_single_line():
    assert count_text("hello world\n") == Counts(lines=1, words=2, chars=12)


def test_counts_empty_text():
    assert count_text("") == Counts(lines=0, words=0, chars=0)
