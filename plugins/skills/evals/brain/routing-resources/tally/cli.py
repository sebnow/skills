"""Command line entry point."""

import argparse
import sys

from . import invoice, report
from .entries import ParseError, parse_log


def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser(prog="tally")
    sub = parser.add_subparsers(dest="command", required=True)

    p_report = sub.add_parser("report", help="hours per client for each ISO week")
    p_report.add_argument("log")

    p_invoice = sub.add_parser("invoice", help="print an invocie for one client")
    p_invoice.add_argument("log")
    p_invoice.add_argument("--client", required=True)
    p_invoice.add_argument("--rate", type=float, required=True)

    args = parser.parse_args(argv)
    try:
        with open(args.log, encoding="utf-8") as f:
            entries = parse_log(f.read())
    except (OSError, ParseError) as err:
        print(f"tally: {err}", file=sys.stderr)
        return 1

    if args.command == "report":
        sys.stdout.write(report.render(report.weekly_totals(entries)))
    else:
        lines = invoice.invoice_lines(entries, args.client, args.rate)
        sys.stdout.write(invoice.render(lines))
    return 0


if __name__ == "__main__":
    sys.exit(main())
