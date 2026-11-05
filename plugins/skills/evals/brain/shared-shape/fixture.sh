#!/usr/bin/env bash
# Copies the shared tidewright stub into the workspace, then this case's
# files over it: the fetch code, the vessel file and the filing guidance the
# scout's report cites.
set -euo pipefail
here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cp -R "$here/../resources/." .
cp -R "$here/resources/." .
