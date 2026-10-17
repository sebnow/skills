#!/usr/bin/env bash
# Builds the Go fixture as a git repository: one commit with deskbook at 0.3.0,
# then one commit per directory under commits/, each an overlay of changed
# files plus a MESSAGE file holding the commit message.
set -euo pipefail
here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

commit() {
	git add -A
	GIT_AUTHOR_DATE="$1" GIT_COMMITTER_DATE="$1" \
		git -c user.name="Dana Okafor" -c user.email="dana@brightloop.dev" \
		commit -q -F "$2"
}

git init -q -b main
cp -R "$here/resources/." .
printf 'deskbook 0.3.0\n' > "$TMPDIR/base-message"
commit "2026-03-02T10:00:00Z" "$TMPDIR/base-message"

day=3
for dir in "$here"/commits/*/; do
	(cd "$dir" && find . -type f ! -name MESSAGE -exec cp --parents {} "$OLDPWD" \;)
	commit "2026-03-0${day}T10:00:00Z" "$dir/MESSAGE"
	day=$((day + 1))
done

# The eval sandbox will not execute git, so export the commits under review
# where Read can reach them.
git format-patch -3 --stdout > last-three-commits.patch
