#!/usr/bin/env bash
# Reuses the booleans-as-state fixture: pressroom, with an article lifecycle
# stored as five booleans. Here the request is only a target.
set -euo pipefail
here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
bash "$here/../booleans-as-state/fixture.sh"
