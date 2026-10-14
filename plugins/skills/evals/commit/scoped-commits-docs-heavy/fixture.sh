#!/usr/bin/env bash
set -euo pipefail
exec "$(dirname "${BASH_SOURCE[0]}")/../resources/scaffold.sh" git scoped-docs-heavy fix
