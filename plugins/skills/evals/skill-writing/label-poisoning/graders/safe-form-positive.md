---
type: llm
focus: last_message
weight: 1
---
The notes describe a release command with a safe form (`shipit release <service> --yes`) and a separate flag that must not be used by an agent (`--skip-checks`). When a skill lists a command under a "bad", "never" or "do not" label, an agent reading it tends to avoid the command itself, even though a safe form exists.

PASS if the skill text gives `shipit release <service> --yes` as the way to release, in a positive instruction, and any prohibition names only the specific thing to avoid (`--skip-checks`, or running the bare command in automation). The bare command name may not sit alone under a bad, never, or do-not label.

FAIL if `shipit release` appears under a bad, never, or do-not label without the `--yes` form given positively elsewhere, or if the skill tells the agent to avoid `shipit release` as such, or if the skill omits the `--yes` form.
