from wordstat.count import Counts, count


def test_counts_a_single_line():
    assert count("hello world\n") == Counts(lines=1, words=2, chars=12)


def test_counts_empty_text():
    assert count("") == Counts(lines=0, words=0, chars=0)


def test_counts_words_separated_by_newlines_and_tabs():
    assert count("one\ntwo\tthree\n").words == 3
