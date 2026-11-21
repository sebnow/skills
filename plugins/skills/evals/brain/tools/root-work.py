"""Counts the tool calls the brain made itself in each run of a
`claude plugin eval --keep-temp --json` result, so a routing case can be
read by who did the work rather than by whether a worker started.

Usage: python3 -I root-work.py <results.json>...

Prints one line per run: case, run, the workers dispatched in order, and
the root's own calls split at the first dispatch, as "before | after".
Root calls are those whose message has no parent_tool_use_id; Agent,
SendMessage, Skill, TodoWrite and the Task tools are coordination and are
not counted as work.

The last column is the workload verdict for the route-* cases other than
route-tiny: no root Edit or Write, at most two Read, Grep and Glob calls
together, and no Bash in route-question or at most one Bash elsewhere
(to check a worker's result).
"""
import json
import sys
from collections import Counter

COORDINATION = {"Agent", "Task", "SendMessage", "Skill", "TodoWrite",
                "TaskCreate", "TaskUpdate", "TaskList", "TaskGet", "TaskOutput",
                "TaskStop"}


def summarise(trace):
    workers, before, after = [], Counter(), Counter()
    for line in open(trace):
        try:
            d = json.loads(line)
        except ValueError:
            continue
        if d.get("type") != "assistant" or d.get("parent_tool_use_id"):
            continue
        for b in d["message"].get("content", []):
            if b.get("type") != "tool_use":
                continue
            name = b["name"]
            if name in ("Agent", "Task"):
                workers.append(b["input"].get("subagent_type") or "general")
            elif name == "SendMessage":
                workers.append("resume")
            elif name not in COORDINATION:
                (after if workers else before)[name] += 1
    return workers, before, after


def verdict(case, c):
    if case == "route-tiny" or not case.startswith("route-"):
        return "n/a"
    edits = c["Edit"] + c["Write"] + c["NotebookEdit"]
    reads = c["Read"] + c["Grep"] + c["Glob"]
    bash = c["Bash"] <= (0 if case == "route-question" else 1)
    return "ok" if edits == 0 and reads <= 2 and bash else "FAIL"


def fmt(c):
    return " ".join(f"{k}={v}" for k, v in sorted(c.items())) or "-"


def main():
    for path in sys.argv[1:]:
        res = json.load(open(path))
        for case in res["cases"]:
            for arm, runs in case.get("arms", {}).items():
                for i, run in enumerate(runs):
                    trace = run.get("tracePath")
                    if not trace:
                        print(f"{case['name']}\t{arm}#{i}\terror: {run.get('error')}")
                        continue
                    workers, before, after = summarise(trace)
                    print(f"{case['name']}\t{arm}#{i}\t"
                          f"workers: {','.join(w.removeprefix('sebnow:') for w in workers) or '-'}\t"
                          f"root: {fmt(before)} | {fmt(after)}\t"
                          f"workload: {verdict(case['name'], before + after)}")


main()
