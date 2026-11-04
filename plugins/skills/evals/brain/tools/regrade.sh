#!/usr/bin/env bash
# Regrades the runs of an earlier `claude plugin eval --keep-temp --json`
# with the graders in a directory, by passing each run's last message
# through probe-grader.sh. Only graders on last_message give the same
# input this way; a trace grader would see one message, not the trace.
# Set JUDGE to choose the judge model.
#
# Usage: regrade.sh <graders dir> <results.json>...
# Prints one line per run and grader: <run>-<grader>, the judge votes, and
# the lengths the judge saw and the reply had.
set -euo pipefail
here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
graders="$(realpath "$1")"; shift
out="$(mktemp -d)"
for results in "$@"; do
  for trace in $(jq -r '.cases[].arms.with[] | select(.error == null) | .tracePath' "$results"); do
    run="$(basename "$(dirname "$(dirname "$trace")")")"; run="${run#claude-eval-}"
    python3 -I "$here/trace-replies.py" "$trace" --last > "$out/$run.txt"
    for g in "$graders"/*.md; do
      "$here/probe-grader.sh" "$g" "$run-$(basename "$g" .md)=$out/$run.txt" 2>/dev/null &
    done
    wait
  done
done
echo "replies: $out" >&2
