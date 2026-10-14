#!/usr/bin/env bash
# Builds a wordstat history in the current git repository. Each line of
# <history-file> is "<step><TAB><subject>"; the steps below run in that
# order and each is committed with its subject. History files that list
# the same steps in the same order build the same trees, so only the
# commit style differs between them.
# Usage: history.sh <history-file>
set -euo pipefail
history="$1"
src="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/wordstat"

put() {
  for f in "$@"; do
    mkdir -p "$(dirname "$f")"
    cp "$src/before/$f" "$f" 2>/dev/null || cp "$src/$f" "$f"
  done
}

skeleton() { put pyproject.toml wordstat/__init__.py .gitignore; }
count() { put wordstat/count.py tests/test_count.py; }
util() { put wordstat/util.py tests/test_util.py; }
cli() {
  put wordstat/cli.py tests/test_cli.py
  sed -i 's/^        prog="wordstat", description="Count lines, words and characters."$/        prog="wordstat"/' wordstat/cli.py
}
readme() { put README.md; }
output() { put docs/output.md; }
treewide() {
  put LICENSE
  for f in $(git ls-files '*.py'); do
    { echo "# SPDX-License-Identifier: MIT"; [ -s "$f" ] && echo; cat "$f"; } >"$f.tmp"
    mv "$f.tmp" "$f"
  done
}
cli_help() {
  sed -i 's/^        prog="wordstat"$/        prog="wordstat", description="Count lines, words and characters."/' wordstat/cli.py
}
install() { put docs/install.md; }
examples() { put docs/examples.md; }
exit_status() {
  cat >>docs/output.md <<'MD'

## Exit status

wordstat exits 0 when every file was read, and 2 when the arguments are
invalid. A file that cannot be read stops the run with a traceback.
MD
}
gen_formats() {
  mkdir -p data tools
  cp "$src/generated/base/data/formats.csv" data/
  cp "$src/generated/base/tools/gen_formats.py" tools/
  cp "$src/generated/base/wordstat/formats.py" wordstat/
}
format_flag() { cp "$src/generated/base/wordstat/cli.py" wordstat/; }
api_docs() {
  cat >>README.md <<'MD'

## Python API

`wordstat.count.count(text)` returns a `Counts` with the `lines`, `words`
and `chars` of a string:

```python
from wordstat.count import count

count("one two\n")  # Counts(lines=1, words=2, chars=8)
```
MD
}
readme_links() {
  cat >>README.md <<'MD'

See [docs/](docs/) for installation, the output format and examples.
MD
}

while IFS=$'\t' read -r step subject; do
  "$step"
  git add -A
  git commit -q -m "$subject"
done <"$history"
