#!/usr/bin/env bash
# Rewrites execution.model and execution.append_system_prompt in the junior
# cases from agents/junior.md, so each case runs the model and body the agent
# ships with. Each entry in `cases` is "<case>:<source>", where
# <source> is one of:
#   (empty)       the junior.md body as it is on disk;
#   rule:<text>   that body without the top-level bullet whose first line
#                 starts with "- <text>";
#   report        that body without the "Report in this shape" paragraph
#                 and the list under it.
# With --check, exits 1 if any case.yaml is out of date and changes nothing.
set -euo pipefail
here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
junior="$here/../../agents/junior.md"

cases=(
  "scope:"
  "scope-baseline:rule:Do exactly what it says"
  "stale-line:"
  "stale-line-baseline:rule:A line number in the brief"
  "gap:"
  "gap-baseline:rule:If a step cannot be done"
  "verify:"
  "verify-baseline:rule:Run the verification"
  "report-shape:"
  "report-shape-baseline:report"
)

# Strips the frontmatter and leading blank lines from junior.md on stdin.
body() { awk 'BEGIN{n=0} /^---$/ && n<2 {n++; next} n>=2' | sed '/./,$!d'; }
trim_trailing_blank() { sed -e :a -e '/^\n*$/{$d;N;ba' -e '}'; }
without_rule() {
  PREFIX="- $1" awk '
    index($0, ENVIRON["PREFIX"]) == 1 { skip = 1; next }
    skip && /^  / { next }
    { skip = 0; print }
  '
}
without_report() { awk '/^Report in this shape/ { exit } { print }' | trim_trailing_blank; }

# Sets execution.model to the given model.
set_model() { MODEL="$1" awk '/^  model: / { print "  model: " ENVIRON["MODEL"]; next } { print }'; }

# Replaces the block under "  append_system_prompt: |" up to the next
# execution-level key with the given text, indented four spaces.
render() {
  local yaml="$1" text="$2"
  TEXT="$text" awk '
    /^  append_system_prompt: \|$/ {
      print; n = split(ENVIRON["TEXT"], lines, "\n")
      for (i = 1; i <= n; i++) print (lines[i] == "" ? "" : "    " lines[i])
      skip = 1; next
    }
    skip && /^  [a-z_]+:/ { skip = 0 }
    !skip
  ' "$yaml"
}

full="$(body < "$junior")"
model="$(awk '/^---$/ { n++; next } n == 1 && /^model: / { print $2 }' "$junior")"
status=0
for entry in "${cases[@]}"; do
  name="${entry%%:*}"; source="${entry#*:}"
  case "$source" in
    "") text="$full" ;;
    report)
      text="$(printf '%s\n' "$full" | without_report)"
      ;;
    rule:*)
      text="$(printf '%s\n' "$full" | without_rule "${source#rule:}")"
      ;;
  esac
  if [[ "$text" == "$full" ]]; then
    [[ -z "$source" ]] || { echo "$name: junior.md has no part matching \"$source\"" >&2; exit 1; }
  fi
  yaml="$here/$name/case.yaml"
  want="$(render "$yaml" "$text" | set_model "$model")"
  if [[ "${1:-}" == "--check" ]]; then
    if ! diff -u "$yaml" <(printf '%s\n' "$want") >/dev/null; then
      echo "$name/case.yaml is out of date; run $0" >&2
      status=1
    fi
  else
    printf '%s\n' "$want" > "$yaml"
  fi
done
exit "$status"
