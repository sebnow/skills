#!/usr/bin/env bash
# Reuses the wrong-home fixture: a Go repository whose last three commits are
# exported to last-three-commits.patch. Here they are the subject of a plain
# code review, not a design review.
set -euo pipefail
here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
bash "$here/../wrong-home/fixture.sh"
