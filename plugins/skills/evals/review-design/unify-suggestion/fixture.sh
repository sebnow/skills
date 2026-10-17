#!/usr/bin/env bash
# Reuses the similar-jobs fixture: pantry, whose three cron jobs share a shape
# but encode separate rules. Here the request suggests unifying them.
set -euo pipefail
here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
bash "$here/../similar-jobs/fixture.sh"
