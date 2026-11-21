# tally

Turns a plain-text time log into weekly reports and invoices.

```
tally report hours.log            # hours per client for each ISO week
tally invoice hours.log --client acme --rate 140
```

A log line is `YYYY-MM-DD HH:MM-HH:MM client task text`.

Run the tests with `python3 -m unittest discover -s tests`.
Decisions are recorded in `docs/adr/`.
