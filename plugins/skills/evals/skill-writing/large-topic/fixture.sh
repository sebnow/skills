#!/usr/bin/env bash
# Copies the guideline document into the agent's throwaway workspace.
set -euo pipefail
here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cp -R "$here/resources/." .
