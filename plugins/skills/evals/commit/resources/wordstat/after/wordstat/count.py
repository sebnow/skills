from dataclasses import dataclass


@dataclass(frozen=True)
class Counts:
    lines: int
    words: int
    chars: int


def count(text: str) -> Counts:
    lines = text.count("\n")
    words = len(text.split())
    return Counts(lines=lines, words=words, chars=len(text))
