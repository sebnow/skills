#!/usr/bin/env bash
# Builds a wordstat repository in the current (empty) workspace and leaves
# changes from resources/wordstat/after uncommitted.
#
# Usage: scaffold.sh git|jj [<history> <change>]
#
# Without <history>: one commit, "Add wordstat CLI", holding
# resources/wordstat/before, and all of resources/wordstat/after
# uncommitted: a word-count bug fix, an unused human_size helper, and a
# --json flag.
#
# With <history> (git only): the commits history.sh builds from
# history/<history>.txt. <change> names what is left
# uncommitted: "all" as above, "fix" for the word-count bug fix alone, or
# "generated" for resources/wordstat/generated/after: a new row in
# data/formats.csv, wordstat/formats.py regenerated from it, and a JSON
# output format in the CLI that uses the new enum member, or "shared" for
# resources/wordstat/shared: a rename of count() to count_text(), a fix
# that reports unreadable files and keeps counting, and a --json flag,
# the three meeting in wordstat/cli.py and tests/test_cli.py.
set -euo pipefail
vcs="$1" history="${2:-}" change="${3:-all}"
here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ws="$PWD"

# The run's sandbox cannot execute /usr/bin/git, and ~/.nix-profile is
# unreadable, but $HOME/.nix-profile/bin is on the agent's PATH and $HOME is
# the run's temporary home. Link the real binaries there.
mkdir -p "$HOME/.nix-profile/bin"
for bin in git jj; do
  ln -sf "$(readlink -f "$(command -v "$bin")")" "$HOME/.nix-profile/bin/$bin"
done

exclude_sandbox_files() {
  # The sandbox mounts these into the workspace; a real repository would
  # not show them as untracked.
  printf '%s\n' .bash_profile .bashrc .claude/ .eval-artifacts .gitconfig \
    .gitmodules .idea .mcp.json .profile .ripgreprc .vscode .zprofile .zshrc \
    >>.git/info/exclude
}

case "$vcs" in
git)
  git init -q -b main
  git config user.name "Sam Doe"
  git config user.email sam@example.com
  exclude_sandbox_files
  if [ -n "$history" ]; then
    "$here/history.sh" "$here/history/$history.txt"
  else
    cp -R "$here/wordstat/before/." .
    git add -A
    git commit -q -m "Add wordstat CLI"
  fi
  state=.git/eval-state.txt
  ;;
jj)
  if [ -n "$history" ]; then
    echo "scaffold.sh: <history> needs git" >&2
    exit 2
  fi
  # Colocated: in a non-colocated repository git finds the sandbox's own
  # repository in the run's home directory, and agents committed the whole
  # project there.
  cp -R "$here/wordstat/before/." .
  JJ_USER="Sam Doe" JJ_EMAIL=sam@example.com jj git init --colocate . >/dev/null 2>&1
  jj config set --repo user.name "Sam Doe" 2>/dev/null
  jj config set --repo user.email sam@example.com 2>/dev/null
  exclude_sandbox_files
  jj commit -m "Add wordstat CLI" >/dev/null 2>&1
  state=.jj/eval-state.txt
  ;;
*)
  echo "usage: scaffold.sh git|jj [<history> <change>]" >&2
  exit 2
  ;;
esac
base="$(git rev-parse HEAD)"

case "$change" in
all) cp -R "$here/wordstat/after/." . ;;
fix)
  # A patch rather than a copy, so the change applies on top of history
  # that edited the same files.
  for f in wordstat/count.py tests/test_count.py; do
    { diff -u "$here/wordstat/before/$f" "$here/wordstat/after/$f" || true; } | patch -s --no-backup-if-mismatch "$f"
  done
  ;;
generated) cp -R "$here/wordstat/generated/after/." . ;;
shared) cp -R "$here/wordstat/shared/." . ;;
*)
  echo "scaffold.sh: unknown change: $change" >&2
  exit 2
  ;;
esac

# Graders read repository state from $state, which this detached poller
# rewrites every second until the workspace is gone. It runs outside the
# agent's sandbox and only reads the repository.
setsid "$here/dump-state.sh" "$vcs" "$ws" "$base" "$state" </dev/null >/dev/null 2>&1 &
