# wordstat

Counts lines, words and characters in text files.

```sh
wordstat notes.txt todo.txt
```

Each file gets one line: lines, words, characters, then the path.

Pass `--format json` to print one JSON object per file instead, for
scripts:

```sh
wordstat --format json notes.txt
{"path": "notes.txt", "lines": 12, "words": 80, "chars": 431}
```
