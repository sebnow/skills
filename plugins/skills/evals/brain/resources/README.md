# tidewright

Turns a coastal craft's voyage logs into the CCVR 2021 operating return.

```
tidewright import logs/2025-season.csv
tidewright return 2025 > return-2025.md
```

The return opens with its status. Until every input a rule needs is recorded
it says **NOT FOR SUBMISSION** and lists the open findings.

## Elections

Where a clause of the Regulations leaves a choice to the operator, record it:

```
tidewright clauses                      # clauses that offer a choice
tidewright elect survey-deferral three-months -evidence surveyor-letter.pdf
tidewright elections                    # in force first, then superseded
```

Electing again on a clause supersedes the earlier election; both stay on
record (REQ-0139, REQ-0140).

See `docs/` for the glossary, requirements, ADRs and the regulation text.
