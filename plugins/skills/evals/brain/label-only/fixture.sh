#!/usr/bin/env bash
# Copies the shared tidewright stub into the workspace, then this case's
# files over it: the overhaul and survey requirements the reports cite, the
# class-scheme notice, and the files the senior's report says it wrote.
set -euo pipefail
here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cp -R "$here/../resources/." .
cp -R "$here/resources/." .
