#!/usr/bin/env bash
# Grades saved replies with one grader through claude plugin eval, without
# running the brain: each reply goes in a throwaway case whose prompt asks
# a small model to echo it, so the grader sees it as the last message.
# Use it to check a rubric on known-good and known-bad replies before
# trusting it; an llm judge in the harness has failed replies that the same
# model passed when given the same rubric directly.
#
# Usage: probe-grader.sh <grader.md> <name>=<reply file> ...
# Prints one line per reply: name, verdict, and the length the judge saw
# against the reply's length (a mismatch means the echo was not verbatim).
set -euo pipefail
grader="$(realpath "$1")"; shift
work="$(mktemp -d)"
mkdir -p "$work/.claude-plugin" "$work/evals"
echo '{"name":"grader-probe","version":"0.0.1"}' > "$work/.claude-plugin/plugin.json"
for spec in "$@"; do
  name="${spec%%=*}"; file="$(realpath "${spec#*=}")"
  mkdir -p "$work/evals/$name/graders"
  cp "$grader" "$work/evals/$name/graders/"
  python3 -I - "$work/evals/$name/case.yaml" "$name" "$file" <<'PY'
import json, sys
out, name, reply = sys.argv[1:]
prompt = ("Reply with exactly the text between the markers, verbatim, and "
          "nothing else. Do not use tools.\n<<<BEGIN\n" + open(reply).read() + "\nEND>>>\n")
json.dump({"schema_version": "1.1", "name": name, "runs": 1,
           "execution": {"max_turns": 1, "timeout_seconds": 300,
                         "allowed_tools": [], "prompt": prompt}},
          open(out, "w"))
PY
  (cd "$work" && claude plugin eval . --case "$name" --ablation none --runs 1 \
    --model haiku --judge-model claude-sonnet-5-5 --no-publish --trust-plugin \
    -j 1 --json "$work/$name.json" >/dev/null 2>&1) &
done
wait
for spec in "$@"; do
  name="${spec%%=*}"; file="${spec#*=}"
  jq -r --arg n "$name" --arg want "$(wc -c < "$file")" \
    '.cases[].arms.with[] | $n + " " + ([.graders[] | .explanation + " saw=" + (.evidence | length | tostring)] | join(" ")) + " reply=" + $want' \
    "$work/$name.json"
done
echo "probe cases and results: $work" >&2
