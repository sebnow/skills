#!/usr/bin/env bash
# Builds the wordstat repository in the current (empty) workspace: one commit
# holding resources/wordstat/before, then the files from
# resources/wordstat/after left uncommitted. Those hold three independent
# changes: a word-count bug fix, an unused human_size helper, and a --json
# flag. Usage: scaffold.sh git|jj
set -euo pipefail
vcs="$1"
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

cp -R "$here/wordstat/before/." .
case "$vcs" in
git)
  git init -q -b main
  git config user.name "Sam Doe"
  git config user.email sam@example.com
  exclude_sandbox_files
  git add -A
  git commit -q -m "Add wordstat CLI"
  state=.git/eval-state.txt
  ;;
jj)
  # Colocated: in a non-colocated repository git finds the sandbox's own
  # repository in the run's home directory, and agents committed the whole
  # project there.
  JJ_USER="Sam Doe" JJ_EMAIL=sam@example.com jj git init --colocate . >/dev/null 2>&1
  jj config set --repo user.name "Sam Doe" 2>/dev/null
  jj config set --repo user.email sam@example.com 2>/dev/null
  exclude_sandbox_files
  jj commit -m "Add wordstat CLI" >/dev/null 2>&1
  state=.jj/eval-state.txt
  ;;
*)
  echo "usage: scaffold.sh git|jj" >&2
  exit 2
  ;;
esac
base="$(git rev-parse HEAD)"
cp -R "$here/wordstat/after/." .

# Graders read repository state from $state, which this detached poller
# rewrites every second until the workspace is gone. It runs outside the
# agent's sandbox and only reads the repository.
setsid "$here/dump-state.sh" "$vcs" "$ws" "$base" "$state" </dev/null >/dev/null 2>&1 &
