#!/usr/bin/env python3
"""Install the packaged brain profile and workers; Python 3.11+, no dependencies."""

import argparse
import hashlib
import json
import os
from pathlib import Path
import tomllib


def digest(data):
    return hashlib.sha256(data).hexdigest()


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--home", type=Path,
                        default=Path(os.environ.get("CODEX_HOME", "~/.codex")).expanduser())
    args = parser.parse_args()
    home = args.home.expanduser().resolve()
    package = Path(__file__).resolve().parents[1]
    profile = (package / "profiles/brain.toml").read_text(encoding="utf-8")
    profile = profile.replace("__BRAIN_SKILL_PATH__", json.dumps(
        str(home / "sebnow/brain/SKILL.md"), ensure_ascii=False)[1:-1])
    files = {"brain.config.toml": profile.encode("utf-8")}
    for name in ("senior", "mid", "junior", "scout"):
        files[f"sebnow/agents/{name}.toml"] = (package / f"agents/{name}.toml").read_bytes()
    for data in files.values():
        tomllib.loads(data.decode("utf-8"))
    files["sebnow/brain/SKILL.md"] = (package / "skills/brain/SKILL.md").read_bytes()

    receipt = home / "sebnow/installed.json"
    previous = json.loads(receipt.read_text()) if receipt.exists() else {}
    # Preflight every destination before writing anything. Never overwrite a
    # file changed by its user, or a file owned by some other installer.
    for relative, data in files.items():
        path = home / relative
        if path.is_symlink():
            raise SystemExit(f"Refusing to replace symlink: {path}")
        if path.exists():
            existing = path.read_bytes()
            if existing != data and previous.get(relative) != digest(existing):
                raise SystemExit(f"Modified or unrelated file; left unchanged: {path}")
    for relative, data in files.items():
        path = home / relative
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_bytes(data)
    receipt.write_text(json.dumps({key: digest(data) for key, data in files.items()}, indent=2) + "\n")
    print(f"Installed brain profile and four workers in {home}")
    print("Start a new session: codex --profile brain")


if __name__ == "__main__":
    main()
