"""Prints what the brain said in an eval run, from the trace.jsonl that
`claude plugin eval --keep-temp` leaves in /tmp/claude-eval-*/out/.

Usage: python3 -I trace-replies.py <trace.jsonl> [--last]

Without --last, prints every top-level text block in order, with markers
for worker dispatches, task notifications and the end of each turn, so the
first reply (everything before the first notification) can be read
whole; an llm judge on `focus: trace` sees only the first and last 12
messages. With --last, prints the run's final message, which is what a
`last_message` grader sees.
"""
import json
import sys


def main():
    path = sys.argv[1]
    last_only = "--last" in sys.argv[2:]
    last = None
    for line in open(path):
        try:
            d = json.loads(line)
        except ValueError:
            continue
        t = d.get("type")
        if t == "result":
            if d.get("result"):
                last = d["result"]
            if not last_only:
                print(f"[end of turn: {d.get('subtype')}, {d.get('num_turns')} turns]")
            continue
        if last_only or d.get("parent_tool_use_id"):
            continue
        if t == "assistant":
            for b in d["message"].get("content", []):
                if b.get("type") == "text":
                    print(b["text"])
                    print("-----")
                elif b.get("type") == "tool_use" and b["name"] in ("Agent", "Task", "SendMessage"):
                    i = b["input"]
                    print(f"[{b['name']}: {i.get('subagent_type')} / {i.get('description')}]")
        elif t == "user":
            c = d["message"].get("content")
            s = c if isinstance(c, str) else " ".join(
                x.get("text", "") for x in c if isinstance(x, dict) and x.get("type") == "text")
            if "task-notification" in s:
                print("[task notification: a worker reported; the first reply ends above]")
    if last_only:
        print(last if last is not None else "")


main()
