from dataclasses import dataclass


@dataclass(frozen=True)
class Counts:
    lines: int
    words: int
    chars: int


def count_text(text: str) -> Counts:
    lines = text.count("\n")
    words = len([w for w in text.split(" ") if w])
    return Counts(lines=lines, words=words, chars=len(text))
