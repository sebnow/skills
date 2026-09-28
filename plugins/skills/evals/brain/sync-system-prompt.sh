#!/usr/bin/env bash
# Rewrites execution.append_system_prompt in the brain cases from
# agents/brain.md. Each entry in `cases` is "<case>:<source>", where
# <source> is one of:
#   (empty)       the brain.md body as it is on disk;
#   <section>     that body without the "# <section>" heading and its text;
#   rev:<rev>     the brain.md body at jujutsu revision <rev>, so a baseline
#                 arm can be the prompt as it was before a change.
# With --check, exits 1 if any case.yaml is out of date and changes nothing.
set -euo pipefail
here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
brain="$here/../../agents/brain.md"

cases=(
  "buried-ask:"
  "buried-ask-baseline:How you report"
  "label-only:"
)

# Strips the frontmatter and leading blank lines from brain.md on stdin.
body() { awk 'BEGIN{n=0} /^---$/ && n<2 {n++; next} n>=2' | sed '/./,$!d'; }
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

full="$(body < "$brain")"
status=0
for entry in "${cases[@]}"; do
  name="${entry%%:*}"; source="${entry#*:}"
  case "$source" in
    "") text="$full" ;;
    rev:*)
      rev="${source#rev:}"
      text="$(cd "$here" && jj --ignore-working-copy file show -r "$rev" "$brain" | body)"
      if [[ -z "$text" ]]; then
        echo "$name: no brain.md body at revision $rev" >&2
        exit 1
      fi
      ;;
    *)
      if ! grep -qxF "# $source" <<<"$full"; then
        echo "$name: brain.md has no \"# $source\" section" >&2
        exit 1
      fi
      text="$(printf '%s\n' "$full" | without_section "$source")"
      ;;
  esac
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
