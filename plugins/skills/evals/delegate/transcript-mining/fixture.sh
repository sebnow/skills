#!/usr/bin/env bash
# Drops a small Claude Code transcript into the agent's throwaway workspace.
set -euo pipefail
here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cp -R "$here/resources/." .
