---
type: llm
focus: last_message
weight: 1
---
finops is a Go library. `internal/rates` loads exchange rates; `Book.Quote(from, to)` returns the latest quote for a pair or an error wrapping `rates.ErrNoRate` (`no rate for currency pair FROM/TO`). Three callers each carry an identical conversion function: `convert` in `invoicing.Total`, `toPartnerCurrency` in `payouts.Amount`, and `inReportingCurrency` in `treasury.Position`. The decision belongs with the rates package (or one conversion function the three share), since it is about how a rate is applied.

What each caller produces today, from the code; every caller behaves the same apart from its error prefix:

| case | result |
|---|---|
| amount already in the target currency | the amount unchanged, no quote looked up |
| no quote for the pair | error wrapping `rates.ErrNoRate` |
| quote's `AsOf` more than 24 hours before now (exactly 24 hours is accepted) | error `FROM/TO rate from <AsOf in RFC 3339> is stale` |
| otherwise | cents × RateMicro / 1,000,000, rounded half to even to the cent, computed without overflow; a result outside int64 is an error `N FROM cents do not fit in TO cents` |
| error prefix | `invoice for <customer>: ` (invoicing), `payout to <partner>: ` (payouts), `cash position: ` (treasury) |

Check each design against this table, not against the report's own description of the code.

The report may offer several designs. PASS if at least one offered design puts this conversion decision in one place that all three callers use (for example a `Convert` method on `rates.Book`, a function in `rates`, a converter type `rates` returns, or one shared conversion function), removes the three copies, and keeps every row of the table, including each caller's error prefix and the 24-hour boundary.

FAIL if no offered design does this: for example every design keeps the three copies, moves only part of the decision (only the staleness check, only the rounding) and leaves the rest copied, or changes a row (a different staleness limit or boundary, a different rounding rule, dropping the same-currency shortcut, or changing the error text callers return).
