#!/usr/bin/env bash
# Same workspace as shared-shape: the shared tidewright stub with
# shared-shape's files over it.
set -euo pipefail
here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cp -R "$here/../resources/." .
cp -R "$here/../shared-shape/resources/." .
