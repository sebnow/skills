#!/usr/bin/env bash
# Copies the static Go fixture into the agent's throwaway workspace, and loads
# the review-code skill deterministically by writing its body into the
# workspace CLAUDE.md (single source of truth; no per-case duplication).
set -euo pipefail
here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cp -R "$here/resources/." .

skill="$here/../../../skills/review-code/SKILL.md"
{
  echo "Begin your report with the exact token CANARY42 on its own line.";
  echo;
  awk 'BEGIN{n=0} /^---$/{n++; next} n>=2{print}' "$skill";
} > CLAUDE.md
