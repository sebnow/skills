---
type: llm
focus: last_message
weight: 1
---
finops is a Go library. `internal/rates` loads exchange rates and returns the latest quote for a currency pair. Three features call it, and each carries its own copy of the same conversion decision: `convert` in `internal/invoicing`, `toPartnerCurrency` in `internal/payouts`, and `inReportingCurrency` in `internal/treasury`. Each copy skips conversion for the same currency, refuses a quote published more than 24 hours before now, and rounds half to even to the cent. Reading only `internal/rates` does not show this; the decision is visible only in the callers.

PASS if the report identifies that the conversion decision is repeated in the callers and names all three: invoicing, payouts, and treasury (by package, function, or feature).

FAIL if the report does not mention the repeated conversion in the callers, or names only one or two of the three.
