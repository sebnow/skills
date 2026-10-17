# finops

Money handling shared by the billing service: exchange rates, invoice
totals, partner payouts, and the treasury cash position. The billing service
imports these packages; this repository has no binary of its own.

Rates come from the treasury team's daily export (`rates.Load`). Every
supported currency has two decimal places, so amounts are whole cents.
