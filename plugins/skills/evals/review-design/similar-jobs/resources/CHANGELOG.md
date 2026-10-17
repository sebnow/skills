# Changelog

## 0.5.0

- `pantry-remind-dormant` emails customers who have not ordered for 60 days,
  once per dormant spell. Cron runs it daily.

## 0.4.0

- `pantry-purge-unverified` deletes signups not verified within 7 days. Cron
  runs it hourly.

## 0.3.0

- Adding an item to a cart reserves its stock. A cart left untouched for 2
  hours is closed; `pantry-expire-carts` returns its stock. Cron runs it every
  10 minutes.
- Cron jobs log to stderr and exit non-zero when they fail.

## 0.2.0

- Checkout totals the cart and records the order time on the account.

## 0.1.0

- Accounts, products, carts. An account must verify its email before it can
  open a cart; a signup not verified within 7 days no longer holds its address.
- Run `pantry-migrate` after each deploy.
