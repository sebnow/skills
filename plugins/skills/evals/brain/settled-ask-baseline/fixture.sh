#!/usr/bin/env bash
# Copies the shared tidewright stub into the workspace, then the files the
# settled-ask cases share over it: the conventions, ADR-0001 and ADR-0004
# reworded to match them, the files the senior's report says it wrote, and
# the session transcript from before compaction under .transcripts/.
set -euo pipefail
here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cp -R "$here/../resources/." .
cp -R "$here/../settled-ask-resources/." .
