#!/usr/bin/env bash
# Copies the research report into the agent's throwaway workspace and creates
# the empty skill directory the agent writes into.
set -euo pipefail
here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cp -R "$here/resources/." .
mkdir -p skills/adr-writing
