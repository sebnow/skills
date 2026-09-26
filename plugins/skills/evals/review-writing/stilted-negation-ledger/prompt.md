---
name: stilted-negation-ledger
tags: [review-writing, stilted]
max_turns: 12
allowed_tools: [Skill, Agent]
runs: 1
---

Review this design note before I commit it, and put the final version in your reply.

```markdown
# Double-entry export

The portfolio model records each transaction as a single row: a buy, a sell,
a dividend, or a transfer. Ledger tools such as hledger and Beancount expect
every transaction to balance across two or more postings. A double-entry
export has to generate what the model does not store: the offsetting posting
for each row.

Most rows map directly. A buy debits the position account and credits the
cash account at the same custodian. A dividend credits income and debits cash.

Transfers need pairing. A cash transfer between custodians is recorded as a
withdrawal at one and a deposit at the other, and the two ends of a cash move
share no identifier. The exporter pairs them by amount and date within a
three-day window and posts unpaired rows against `Equity:Unreconciled`.

Positions moved between custodians keep their lots; cash has no lots. The
exporter carries cost basis across for positions and treats cash as fungible.

The export is one-way. Edits made in the ledger file are not read back.
```
