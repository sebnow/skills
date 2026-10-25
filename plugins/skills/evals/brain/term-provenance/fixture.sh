#!/usr/bin/env bash
# Copies the shared tidewright stub into the workspace, so the dedupe code,
# REQ-0108 and the Regulations the scout's report cites exist.
set -euo pipefail
here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cp -R "$here/../resources/." .
