#!/usr/bin/env bash
# Copies the tidewright stub shared by both arms into the workspace, so the
# files the history's scout reports cite exist at the cited lines.
set -euo pipefail
here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cp -R "$here/../resources/." .
