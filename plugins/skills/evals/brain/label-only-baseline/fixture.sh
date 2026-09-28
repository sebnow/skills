#!/usr/bin/env bash
# Same workspace as label-only: the shared tidewright stub with label-only's
# files over it.
set -euo pipefail
here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cp -R "$here/../resources/." .
cp -R "$here/../label-only/resources/." .
