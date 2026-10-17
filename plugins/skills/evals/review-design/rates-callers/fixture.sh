#!/usr/bin/env bash
# Copies the static Go fixture into the agent's throwaway workspace.
set -euo pipefail
here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cp -R "$here/resources/." .
