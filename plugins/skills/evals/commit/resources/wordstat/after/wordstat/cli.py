import argparse
import json
import sys

from wordstat.count import count
from wordstat.util import read_text


def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser(
        prog="wordstat", description="Count lines, words and characters."
    )
    parser.add_argument("paths", nargs="+")
    parser.add_argument(
        "--json", action="store_true", help="print one JSON object per file"
    )
    args = parser.parse_args(argv)

    for path in args.paths:
        c = count(read_text(path))
        if args.json:
            print(
                json.dumps(
                    {"path": path, "lines": c.lines, "words": c.words, "chars": c.chars}
                )
            )
        else:
            print(f"{c.lines:>7} {c.words:>7} {c.chars:>7} {path}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
