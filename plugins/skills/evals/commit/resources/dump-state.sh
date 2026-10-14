#!/usr/bin/env bash
# Usage: dump-state.sh git|jj <workspace> <base-commit-id> <out>
# Writes, for every commit after <base-commit-id> in topological order:
#   SUBJECT: <first line of the description>
#   FILES: <changed paths, sorted, space-separated>
# then one line naming the uncommitted paths:
#   UNCOMMITTED: <paths>
# In a jj repository a working-copy commit without a description counts as
# uncommitted. The jj repository is colocated, so the agent may commit with
# git instead; while git's HEAD is not jj's @-, jj has not imported those
# commits yet and git's view is the current one.
set -uo pipefail
vcs="$1" ws="$2" base="$3" out="$4"

git_view() {
  git -C "$ws" log --branches HEAD --reverse --topo-order "^$base" \
    --format='SUBJECT: %s' --name-only |
    awk '/^SUBJECT: /{if (s) {print s; print "FILES:" f}; s=$0; f=""; next}
         NF{f=f " " $0}
         END{if (s) {print s; print "FILES:" f}}'
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
