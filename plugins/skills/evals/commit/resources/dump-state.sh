#!/usr/bin/env bash
# Usage: dump-state.sh git|jj <workspace> <base-commit-id> <out> [<remote-git-dir>]
# Writes, for every commit after <base-commit-id> in topological order:
#   SUBJECT: <first line of the description>
#   FILES: <changed paths, sorted, space-separated>
#   WORDS: <changed word-diff runs, sorted, " | "-separated>
#   BODY: <one line per line of the message body>
# then one line naming the uncommitted paths:
#   UNCOMMITTED: <paths>
# With <remote-git-dir>, the same lines prefixed R- for the commits on the
# remote's main after <base-commit-id>, then
#   R-BASE: kept|lost
# saying whether <base-commit-id> is still on the remote's main.
# In a jj repository a working-copy commit without a description counts as
# uncommitted. In a colocated jj repository the agent may commit with git
# instead; while git's HEAD is not jj's @-, jj has not imported those
# commits yet and git's view is the current one.
set -uo pipefail
vcs="$1" ws="$2" base="$3" out="$4" remote="${5:-}"

# Prints one commit in the format above, reading it with "${g[@]}", a git
# command for the repository that stores it.
emit() {
  local c="$1"
  "${g[@]}" log -1 --format='SUBJECT: %s' "$c"
  printf 'FILES:'
  "${g[@]}" diff-tree --no-commit-id --name-only -r "$c" | LC_ALL=C sort | sed 's/^/ /' | tr -d '\n'
  echo
  # The changed runs of a word diff, sorted and deduplicated, so a rename
  # shows as nothing but the old and new identifier.
  echo "WORDS:$("${g[@]}" show --format= --word-diff=porcelain \
    --word-diff-regex='[A-Za-z_][A-Za-z0-9_]*|[^[:space:]]' "$c" |
    grep -E '^[-+]' | grep -vE '^(\+\+\+|---) ' | LC_ALL=C sort -u |
    sed 's/^/ /' | paste -sd'|' | sed 's/|/ |/g')"
  "${g[@]}" log -1 --format=%b "$c" | sed 's/^/BODY: /'
}

git_view() {
  local c g=(git -C "$ws")
  for c in $(git -C "$ws" rev-list --reverse --topo-order --branches HEAD "^$base"); do
    emit "$c"
  done
  printf 'UNCOMMITTED:'
  git -C "$ws" --no-optional-locks status --porcelain --untracked-files=all |
    cut -c4- | sort | tr '\n' ' ' | sed 's/^/ /; s/ $//'
  echo
}

# jj stores its commits as git commits, in .git when colocated and in
# .jj/repo/store/git otherwise, so emit reads them with git.
jj_view() {
  local id flags g
  if [ -d "$ws/.git" ]; then g=(git -C "$ws"); else g=(git --git-dir="$ws/.jj/repo/store/git"); fi
  jj -R "$ws" --ignore-working-copy --color=never log --no-graph --reversed \
    -r "($base::) ~ $base" \
    -T 'commit_id ++ " " ++ if(description, "d", "-") ++ if(empty, "e", "-") ++ "\n"' |
    while read -r id flags; do
      case "$flags" in
      d?) emit "$id" ;;
      -e) ;;
      *)
        printf 'UNCOMMITTED:'
        "${g[@]}" diff-tree --no-commit-id --name-only -r "$id" | LC_ALL=C sort | sed 's/^/ /' | tr -d '\n'
        echo
        ;;
      esac
    done
}

remote_view() {
  local c g=(git --git-dir="$remote")
  for c in $(git --git-dir="$remote" rev-list --reverse --topo-order main "^$base"); do
    emit "$c"
  done | sed 's/^/R-/'
  if git --git-dir="$remote" merge-base --is-ancestor "$base" main 2>/dev/null; then
    echo "R-BASE: kept"
  else
    echo "R-BASE: lost"
  fi
}

snapshot() {
  [ -z "$remote" ] || remote_view
  case "$vcs" in
  git) git_view ;;
  jj)
    if [ ! -d "$ws/.git" ]; then
      jj_view
      return
    fi
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
