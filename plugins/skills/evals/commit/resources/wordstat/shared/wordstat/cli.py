import argparse
import json
import sys

from wordstat.count import count_text
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

    status = 0
    for path in args.paths:
        try:
            text = read_text(path)
        except OSError as e:
            print(f"wordstat: {path}: {e.strerror}", file=sys.stderr)
            status = 1
            continue
        c = count_text(text)
        if args.json:
            print(
                json.dumps(
                    {"path": path, "lines": c.lines, "words": c.words, "chars": c.chars}
                )
            )
        else:
            print(f"{c.lines:>7} {c.words:>7} {c.chars:>7} {path}")
    return status


if __name__ == "__main__":
    sys.exit(main())
