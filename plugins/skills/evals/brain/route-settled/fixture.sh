#!/usr/bin/env bash
# Copies the tally project shared by the routing cases into the workspace
# and commits it, so a worker can run the tests and commit on top.
set -euo pipefail
here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cp -R "$here/../routing-resources/." .
git init -q -b main
git config user.name "Tally Dev"
git config user.email "dev@example.com"
git add -A
git commit -q -m "tally 1.3.0 with unreleased changes"
