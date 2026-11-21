import unittest
from datetime import date

from tally.entries import ParseError, parse_line, parse_log


class ParseLineTest(unittest.TestCase):
    def test_fields(self):
        e = parse_line("2026-03-02 09:00-10:30 acme sprint planning")
        self.assertEqual(e.day, date(2026, 3, 2))
        self.assertEqual(e.client, "acme")
        self.assertEqual(e.task, "sprint planning")
        self.assertEqual(e.minutes, 90)

    def test_end_before_start(self):
        with self.assertRaises(ParseError):
            parse_line("2026-03-02 10:00-09:00 acme oops")

    def test_skips_comments_and_blanks(self):
        log = "# week 10\n\n2026-03-02 09:00-09:10 acme call\n"
        self.assertEqual(len(parse_log(log)), 1)


if __name__ == "__main__":
    unittest.main()
