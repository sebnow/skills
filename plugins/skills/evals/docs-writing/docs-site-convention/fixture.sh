#!/usr/bin/env bash
# Copies the inkwell service shared with the architecture-visualization evals,
# then this case's existing documentation site on top.
set -euo pipefail
here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cp -R "$here/../../architecture-visualization/resources/." .
cp -R "$here/resources/." .
