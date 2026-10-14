#!/usr/bin/env bash
# Builds a nine-commit wordstat history in the current git repository, one
# commit per line of <subjects-file>. Every subjects file describes the same
# nine steps, so the trees match and only the commit style differs. The
# history ends with docs: entries on top, the reported trigger for agents
# reading the log as Conventional Commits.
# Usage: history.sh <subjects-file>
set -euo pipefail
subjects="$1"
src="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/wordstat"
step=0

put() {
  for f in "$@"; do
    mkdir -p "$(dirname "$f")"
    cp "$src/before/$f" "$f" 2>/dev/null || cp "$src/$f" "$f"
  done
}

commit() {
  step=$((step + 1))
  git add -A
  git commit -q -m "$(sed -n "${step}p" "$subjects")"
}

put pyproject.toml wordstat/__init__.py .gitignore
commit
put wordstat/count.py tests/test_count.py
commit
put wordstat/util.py tests/test_util.py
commit
put wordstat/cli.py tests/test_cli.py
sed -i 's/^        prog="wordstat", description="Count lines, words and characters."$/        prog="wordstat"/' wordstat/cli.py
commit
put README.md
commit
put docs/output.md
commit
put LICENSE
for f in $(git ls-files '*.py'); do
  { echo "# SPDX-License-Identifier: MIT"; [ -s "$f" ] && echo; cat "$f"; } >"$f.tmp"
  mv "$f.tmp" "$f"
done
commit
sed -i 's/^        prog="wordstat"$/        prog="wordstat", description="Count lines, words and characters."/' wordstat/cli.py
commit
put docs/install.md
commit
[ "$step" -eq "$(wc -l <"$subjects")" ]
