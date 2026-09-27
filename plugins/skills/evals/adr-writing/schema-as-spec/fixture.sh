#!/usr/bin/env bash
# Copies the migration and the design note into the agent's workspace.
set -euo pipefail
here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cp -R "$here/resources/." .
