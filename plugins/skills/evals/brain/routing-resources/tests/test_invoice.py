import unittest

from tally.entries import parse_log
from tally.invoice import invoice_lines, render

LOG = """\
2026-03-02 09:00-09:10 acme call
2026-03-03 13:00-14:20 acme review
2026-03-04 10:00-10:30 globex audit
"""


class InvoiceTest(unittest.TestCase):
    def test_lines_for_one_client(self):
        lines = invoice_lines(parse_log(LOG), "acme", 100)
        self.assertEqual([l.task for l in lines], ["call", "review"])
        self.assertEqual(lines[0].hours, 0.2)

    def test_total(self):
        out = render(invoice_lines(parse_log(LOG), "acme", 100))
        self.assertIn("total", out)


if __name__ == "__main__":
    unittest.main()
