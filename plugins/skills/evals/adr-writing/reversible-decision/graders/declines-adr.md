---
type: llm
focus: last_message
weight: 1
---
The user asked for an ADR for flipping one line in a dev-only config file (`LOG_FORMAT=json` to `text`) for readability. Staging and prod are unaffected. The change is reversible by editing the same line back, and nobody reading the repo later would be surprised by it or need its history explained, so it does not warrant an ADR.

PASS if the reply tells the user this change does not need an ADR, giving as the reason that it is cheap to reverse, local to dev, or unlikely to puzzle a later reader (any of these), and does not produce an ADR. Suggesting a commit message, a comment next to the config line, or just making the change is fine.

FAIL if the reply contains an ADR or ADR-shaped document (a titled record with sections such as Status, Context, Decision, Consequences), including when it first questions whether an ADR is needed and then writes one anyway, or offers a "short version".
