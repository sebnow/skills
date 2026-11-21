import unittest

from tally.entries import parse_log
from tally.report import render, weekly_totals

LOG = """\
2026-03-02 09:00-09:10 acme call
2026-03-03 13:00-14:20 acme review
2026-03-04 10:00-10:30 globex audit
2026-03-09 09:00-10:00 acme planning
"""


class WeeklyTotalsTest(unittest.TestCase):
    def test_rounds_each_entry_up_to_quarter_hour(self):
        totals = weekly_totals(parse_log(LOG))
        self.assertEqual(totals["2026-W10"]["acme"], 15 + 90)
        self.assertEqual(totals["2026-W10"]["globex"], 30)
        self.assertEqual(totals["2026-W11"]["acme"], 60)

    def test_render(self):
        out = render(weekly_totals(parse_log(LOG)))
        self.assertIn("2026-W10", out)
        self.assertIn("1.75 h", out)


if __name__ == "__main__":
    unittest.main()
