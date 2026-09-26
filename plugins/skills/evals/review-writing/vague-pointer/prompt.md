---
name: vague-pointer
tags: [review-writing, pointer]
max_turns: 12
allowed_tools: [Skill, Agent]
runs: 1
---

Review this design note before I commit it, and put the final version in your reply.

```markdown
# Beancount export

The portfolio database records each transaction as a single row. The
`ledger-export` command converts those rows into a Beancount file so that
balances and cost basis can be checked with `bean-check` and `fava`.

## Mapping

A buy debits the position account and credits the cash account at the same
custodian. A sell does the reverse and books the gain or loss to
`Income:CapitalGains`. A dividend credits `Income:Dividends` and debits cash.
Fees are posted to `Expenses:Fees` on the same transaction.

Account names follow the pattern `Assets:<Custodian>:<Symbol>` for positions
and `Assets:<Custodian>:Cash` for cash. Beancount syntax is documented
elsewhere and is not covered here.

## Transfers

A cash transfer between custodians is recorded as a withdrawal at one and a
deposit at the other. The exporter pairs them by amount and date within a
three-day window and posts unpaired rows against `Equity:Unreconciled`.

## Scope

The export is one-way. Edits made in the Beancount file are not read back
into the database.
```
