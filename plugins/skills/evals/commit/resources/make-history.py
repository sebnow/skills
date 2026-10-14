"""Write the resumed transcript for a fix-after-commit case.

The transcript is the session in which the agent committed the word-count
fix: it reads the working copy, commits, and reports. Outputs are read
from a repository built by scaffold.sh with the debug-commit change, whose
ids are deterministic, so the transcript quotes the ids the run will see.

Usage: python3 make-history.py git|jj-only <scaffolded-repo> <out.jsonl>
"""

import json
import subprocess
import sys
import uuid

vcs, repo, out = sys.argv[1:4]
session = str(uuid.UUID(int=1 if vcs == "git" else 2))


def sh(*args: str) -> str:
    return subprocess.run(
        args, cwd=repo, check=True, capture_output=True, text=True
    ).stdout


if vcs == "git":
    status = " M tests/test_count.py\n M wordstat/count.py\n"
    diff = sh("git", "diff", "HEAD~1", "HEAD")
    log = sh("git", "log", "-5", "--format=%s", "HEAD~1")
    look = "git status --short && git diff && git log -5 --format=%s"
    look_out = status + diff + log
    subject = sh("git", "log", "-1", "--format=%s").strip()
    body = sh("git", "log", "-1", "--format=%b").strip()
    short = sh("git", "rev-parse", "--short", "HEAD").strip()
    commit = (
        "git add wordstat/count.py tests/test_count.py && "
        f'git commit -m "{subject}" -m "{body}"'
    )
    stat = sh("git", "show", "--stat", "--format=", "HEAD").strip().splitlines()[-1]
    commit_out = f"[main {short}] {subject}\n {stat.strip()}\n"
    report = (
        f"Committed `{short}` `{subject}`: the fix in `wordstat/count.py` "
        "and its test in `tests/test_count.py`. I didn't run the tests."
    )
else:
    tpl = 'change_id.shortest(8) ++ " " ++ commit_id.shortest(8)'
    ops = sh("jj", "op", "log", "--no-graph", "-T", 'id.short() ++ " " ++ description ++ "\\n"')
    # The operation that committed the agent's change, and the one before.
    lines = [l for l in ops.splitlines() if l.strip()]
    split = sh("jj", "log", "--no-graph", "-r", 'description(glob:"Split words*")', "-T", "change_id")
    at_commit = None
    for line in lines:
        op, _, what = line.partition(" ")
        if not what.startswith("commit"):
            continue
        parent = sh("jj", "--at-op", op, "--ignore-working-copy", "log", "--no-graph", "-r", "@-", "-T", "change_id")
        if parent == split:
            at_commit = op
            break
    before = sh("jj", "--at-op", at_commit, "--ignore-working-copy", "op", "log", "--no-graph", "-n", "2", "-T", 'id.short() ++ "\\n"').split()[1]
    look = "jj st && jj diff --git && jj log -n 5 --no-graph -T 'description.first_line() ++ \"\\n\"'"
    st = sh("jj", "--at-op", before, "--ignore-working-copy", "st")
    diff = sh("jj", "--at-op", before, "--ignore-working-copy", "diff", "--git")
    log = sh("jj", "--at-op", before, "--ignore-working-copy", "log", "-n", "5", "--no-graph", "-T", 'description.first_line() ++ "\\n"')
    look_out = st + diff + log
    desc = sh("jj", "log", "--no-graph", "-r", split, "-T", "description").strip()
    subject, _, body = desc.partition("\n\n")
    commit = f'jj commit -m "{subject}" -m "{body.strip()}"'
    wc = sh("jj", "--at-op", at_commit, "--ignore-working-copy", "log", "--no-graph", "-r", "@", "-T", tpl)
    par = sh("jj", "--at-op", at_commit, "--ignore-working-copy", "log", "--no-graph", "-r", "@-", "-T", tpl)
    commit_out = (
        f"Working copy  (@) now at: {wc} (empty) (no description set)\n"
        f"Parent commit (@-)      : {par} {subject}\n"
    )
    report = (
        f"Committed `{par.split()[0]}` \"{subject}\": the fix in "
        "`wordstat/count.py` and its test in `tests/test_count.py`. "
        "I didn't run the tests."
    )

entries = []
parent = None
clock = 0


def add(kind: str, message: dict) -> None:
    global parent, clock
    clock += 1
    u = str(uuid.UUID(int=(1 if vcs == "git" else 2) * 1000 + clock))
    entry = {
        "parentUuid": parent,
        "isSidechain": False,
        "type": kind,
        "message": message,
        "uuid": u,
        "timestamp": f"2026-10-01T09:0{clock}:00.000Z",
        "userType": "external",
        "entrypoint": "sdk-cli",
        "cwd": "/work",
        "sessionId": session,
        "version": "2.1.293",
    }
    entries.append(entry)
    parent = u


def assistant(content: list) -> None:
    add("assistant", {
        "model": "claude-sonnet-5-5", "id": f"msg_{clock:04d}", "type": "message",
        "role": "assistant", "content": content, "stop_reason": "tool_use"
        if content[-1]["type"] == "tool_use" else "end_turn",
        "stop_sequence": None,
        "usage": {"input_tokens": 1, "output_tokens": 1},
    })


def bash(n: int, command: str, output: str) -> None:
    tid = f"toolu_{vcs.replace('-', '')}{n:02d}"
    assistant([{"type": "tool_use", "id": tid, "name": "Bash",
                "input": {"command": command, "description": "Inspect or commit"}}])
    add("user", {"role": "user", "content": [
        {"tool_use_id": tid, "type": "tool_result", "content": output, "is_error": False}]})


add("user", {"role": "user", "content": "Commit my changes."})
bash(1, look, look_out)
bash(2, commit, commit_out)
assistant([{"type": "text", "text": report}])

with open(out, "w", encoding="utf-8") as f:
    for e in entries:
        f.write(json.dumps(e) + "\n")
