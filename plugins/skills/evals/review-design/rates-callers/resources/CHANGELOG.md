# Changelog

## 0.4.0

- `treasury.Position` values cash balances held in several currencies in one
  reporting currency, at the latest rates.

## 0.3.0

- `payouts.Amount` pays a partner's earnings in the partner's currency.

## 0.2.0

- `invoicing.Total` converts lines priced in another currency into the
  invoice currency. A rate published more than a day before the conversion is
  refused as stale. Converted amounts round half to even, to the cent.

## 0.1.0

- `rates.Load` reads the treasury export; `Book.Quote` returns the latest
  rate for a currency pair.
