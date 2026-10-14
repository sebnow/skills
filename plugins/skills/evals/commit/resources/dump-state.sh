#!/usr/bin/env bash
# Usage: dump-state.sh git|jj <workspace> <base-commit-id> <out>
# Writes, for every commit after <base-commit-id> in topological order:
#   SUBJECT: <first line of the description>
#   FILES: <changed paths, sorted, space-separated>
#   WORDS: <changed word-diff runs, sorted, " | "-separated> (git only)
#   BODY: <one line per line of the message body> (git only)
# then one line naming the uncommitted paths:
#   UNCOMMITTED: <paths>
# In a jj repository a working-copy commit without a description counts as
# uncommitted. The jj repository is colocated, so the agent may commit with
# git instead; while git's HEAD is not jj's @-, jj has not imported those
# commits yet and git's view is the current one.
set -uo pipefail
vcs="$1" ws="$2" base="$3" out="$4"

git_view() {
  local c
  for c in $(git -C "$ws" rev-list --reverse --topo-order --branches HEAD "^$base"); do
    git -C "$ws" log -1 --format='SUBJECT: %s' "$c"
    printf 'FILES:'
    git -C "$ws" diff-tree --no-commit-id --name-only -r "$c" | LC_ALL=C sort | sed 's/^/ /' | tr -d '\n'
    echo
    # The changed runs of a word diff, sorted and deduplicated, so a rename
    # shows as nothing but the old and new identifier.
    echo "WORDS:$(git -C "$ws" show --format= --word-diff=porcelain \
      --word-diff-regex='[A-Za-z_][A-Za-z0-9_]*|[^[:space:]]' "$c" |
      grep -E '^[-+]' | grep -vE '^(\+\+\+|---) ' | LC_ALL=C sort -u |
      sed 's/^/ /' | paste -sd'|' | sed 's/|/ |/g')"
    git -C "$ws" log -1 --format=%b "$c" | sed 's/^/BODY: /'
  done
  printf 'UNCOMMITTED:'
  git -C "$ws" --no-optional-locks status --porcelain --untracked-files=all |
    cut -c4- | sort | tr '\n' ' ' | sed 's/^/ /; s/ $//'
  echo
}

jj_view() {
  jj -R "$ws" --ignore-working-copy --color=never log --no-graph --reversed \
    -r "($base::) ~ $base" \
    -T 'if(description, "SUBJECT: " ++ description.first_line() ++ "\nFILES: " ++ self.diff().files().map(|f| f.path()).join(" ") ++ "\n", if(!empty, "UNCOMMITTED: " ++ self.diff().files().map(|f| f.path()).join(" ") ++ "\n"))'
}

snapshot() {
  case "$vcs" in
  git) git_view ;;
  jj)
    head="$(git -C "$ws" rev-parse -q --verify HEAD)"
    parent="$(jj -R "$ws" --ignore-working-copy log --no-graph -r @- -T commit_id)"
    if [ "$head" = "$parent" ]; then jj_view; else git_view; fi
    ;;
  esac
}

for _ in $(seq 1 1800); do
  [ -d "$ws" ] || exit 0
  if s="$(snapshot 2>&1)"; then
    printf '%s\n' "$s" >"$ws/$out.tmp" 2>/dev/null && mv -f "$ws/$out.tmp" "$ws/$out" 2>/dev/null
  fi
  sleep 1
done
