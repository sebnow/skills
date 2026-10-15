#!/usr/bin/env bash
# Copies the bramble stub shared by the answer-fidelity cases into the
# workspace, so the files the history's lookups and reports cite exist at
# the cited lines.
set -euo pipefail
here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cp -R "$here/../answer-fidelity-resources/." .
