---
name: setup
description: Install or update the sebnow brain profile and worker definitions after installing this Codex plugin.
---

Run `python3 ../../scripts/setup.py`, resolving the script relative to this
skill's directory. It installs the brain profile and workers in CODEX_HOME
(default ~/.codex). Follow the execution environment's permission requirements.
If it reports a modified-file conflict, show that conflict to the user; do not
overwrite it. After success, tell the user to start `codex --profile brain`.
Run setup again after plugin updates to refresh the installed definitions.
