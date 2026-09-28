#!/usr/bin/env bash
# Rewrites execution.append_system_prompt in the brain cases from
# agents/brain.md. Each entry in `cases` is "<case>:<section>": the case gets
# the brain.md body without the "# <section>" heading and its text, or the
# full body when <section> is empty. A baseline arm for a case is one more
# entry naming the section under test. With --check, exits 1 if any case.yaml
# is out of date and changes nothing.
set -euo pipefail
here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
brain="$here/../../agents/brain.md"

cases=(
  "buried-ask:"
  "buried-ask-baseline:How you report"
  "label-only:"
)

body() { awk 'BEGIN{n=0} /^---$/ && n<2 {n++; next} n>=2' "$brain" | sed '/./,$!d'; }
without_section() {
  SECTION="$1" awk '$0 == "# " ENVIRON["SECTION"] {skip=1; next} /^# / {skip=0} !skip' |
    sed -e :a -e '/^\n*$/{$d;N;ba' -e '}'
}

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

full="$(body)"
status=0
for entry in "${cases[@]}"; do
  name="${entry%%:*}"; section="${entry#*:}"
  text="$full"
  if [[ -n "$section" ]]; then
    if ! grep -qxF "# $section" <<<"$full"; then
      echo "$name: brain.md has no \"# $section\" section" >&2
      exit 1
    fi
    text="$(printf '%s\n' "$full" | without_section "$section")"
  fi
  yaml="$here/$name/case.yaml"
  want="$(render "$yaml" "$text")"
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
