#!/usr/bin/env bash
# Copies the ledger-export module and its check.sh into the workspace.
set -euo pipefail
here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cp -R "$here/../resources/." .
