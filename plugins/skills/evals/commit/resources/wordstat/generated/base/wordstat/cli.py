import argparse
import sys

from wordstat.count import count
from wordstat.formats import DESCRIPTIONS, Format
from wordstat.util import read_text


def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser(
        prog="wordstat", description="Count lines, words and characters."
    )
    parser.add_argument("paths", nargs="+")
    parser.add_argument(
        "--format",
        type=Format,
        choices=list(Format),
        default=Format.TEXT,
        metavar="{" + ",".join(f.value for f in Format) + "}",
        help="; ".join(f"{f.value}: {d}" for f, d in DESCRIPTIONS.items()),
    )
    args = parser.parse_args(argv)

    for path in args.paths:
        c = count(read_text(path))
        print(f"{c.lines:>7} {c.words:>7} {c.chars:>7} {path}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
