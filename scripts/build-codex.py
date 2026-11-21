#!/usr/bin/env python3
"""Build the committed Codex package from the Claude plugin. Python 3.11+."""

import json
from pathlib import Path
import tomllib

ROOT = Path(__file__).resolve().parents[1]
SOURCE = ROOT / "plugins/skills"
OUTPUT = ROOT / "plugins/codex"
MODELS = {
    "brain": ("gpt-6.1-sol", "high"),
    "senior": ("gpt-6.1-sol", "medium"),
    "junior": ("gpt-6-luna", "high"),
    "scout": ("gpt-6-luna", "high"),
}


def write(path, text):
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(text, encoding="utf-8")


def toml(values):
    # JSON strings are also valid TOML basic strings for these text values.
    text = "".join(f"{key} = {json.dumps(value, ensure_ascii=False)}\n"
                   for key, value in values.items())
    tomllib.loads(text)
    return text


def main():
    # All files under OUTPUT are generated; use the manifest to remove stale
    # generated files without deleting unexpected files added by a user.
    manifest = OUTPUT / "generated-files.json"
    previous = json.loads(manifest.read_text()) if manifest.exists() else []
    files = {}
    for path in sorted((SOURCE / "skills").rglob("*")):
        if path.is_file():
            files[path.relative_to(SOURCE).as_posix()] = path.read_text(encoding="utf-8")

    review = "skills/review-writing/SKILL.md"
    files[review] = files[review].replace(
        "Launch a subagent (the Agent tool, a fresh general-purpose agent - not a\n   fork, which would inherit this context).",
        "Launch a fresh general-purpose subagent with no inherited conversation\n   (use fork_turns=\"none\").",
    )
    files["skills/skill-writing/SKILL.md"] += (
        "\n## Codex evaluation\n\n"
        "The bundled eval-runner reference describes Claude Code's eval CLI. "
        "Use it only when evaluating the Claude package. For Codex, validate "
        "skills with the available Codex tooling and report which behavioral "
        "checks were actually run.\n"
    )

    for name, (model, effort) in MODELS.items():
        source = (SOURCE / f"agents/{name}.md").read_text(encoding="utf-8")
        _, header, body = source.split("---\n", 2)
        # These source files use a single JSON-quoted line for description.
        description = next(line.removeprefix("description: ")
                           for line in header.splitlines()
                           if line.startswith("description: "))
        description = json.loads(description)
        body = body.strip() + "\n"
        for role in ("senior", "junior", "scout"):
            body = body.replace(f"sebnow:{role}", f"sebnow_{role}")
        body = body.replace("with SendMessage", "with a follow-up message")
        if name == "brain":
            files["skills/brain/SKILL.md"] = (
                "---\nname: brain\ndescription: "
                + json.dumps("Coordinate interactive design and decision sessions; delegate substantial work to the sebnow senior, junior, and scout workers.")
                + "\n---\n\nThese coordinator instructions apply to the primary session. "
                "Delegated workers follow their assigned role.\n\n" + body
            )
            files["profiles/brain.toml"] = toml({
                "model": model,
                "model_reasoning_effort": effort,
                "developer_instructions": (
                    "At session start, read and follow the brain skill at "
                    "__BRAIN_SKILL_PATH__. Use it as the primary session workflow. "
                    "Workers follow their assigned role instead."
                ),
            }) + "\n[agents]\nenabled = true\nmax_concurrent_threads_per_session = 3\n"
        else:
            if name == "scout":
                body += "\nDo not spawn additional agents.\n"
            files[f"agents/{name}.toml"] = toml({
                "model": model,
                "model_reasoning_effort": effort,
                "developer_instructions": body,
            })
            files["profiles/brain.toml"] += (
                f"\n[agents.sebnow_{name}]\n"
                + toml({"description": description,
                        "config_file": f"sebnow/agents/{name}.toml"})
            )

    metadata = json.loads((SOURCE / ".claude-plugin/plugin.json").read_text())
    metadata.update(
        description="Personal agent skills and the brain workflow for Codex.",
        skills="./skills/",
        extensions={"com.openai": {"onboardingSkill": "./skills/setup/SKILL.md"}},
    )
    files[".codex-plugin/plugin.json"] = json.dumps(metadata, indent=2) + "\n"
    files["scripts/setup.py"] = (ROOT / "scripts/setup-codex.py").read_text()
    files["skills/setup/SKILL.md"] = """---
name: setup
description: Install or update the sebnow brain profile and worker definitions after installing this Codex plugin.
---

Run `python3 ../../scripts/setup.py`, resolving the script relative to this
skill's directory. It installs the brain profile and workers in CODEX_HOME
(default ~/.codex). Follow the execution environment's permission requirements.
If it reports a modified-file conflict, show that conflict to the user; do not
overwrite it. After success, tell the user to start `codex --profile brain`.
Run setup again after plugin updates to refresh the installed definitions.
"""
    for relative, content in files.items():
        write(OUTPUT / relative, content)
    for relative in set(previous) - files.keys():
        path = OUTPUT / relative
        if not path.resolve().is_relative_to(OUTPUT.resolve()):
            raise ValueError(f"Invalid generated path: {relative}")
        path.unlink(missing_ok=True)
    write(manifest, json.dumps(sorted(files), indent=2) + "\n")
    tomllib.loads(files["profiles/brain.toml"])
    print(f"Built {len(files)} files in {OUTPUT.relative_to(ROOT)}")


if __name__ == "__main__":
    main()
