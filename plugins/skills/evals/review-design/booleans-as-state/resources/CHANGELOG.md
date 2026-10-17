# Changelog

## 0.4.0

- `POST /articles/{id}/archive` lets the author take a published article off
  the published list. Archiving is final: an archived article cannot be edited, submitted,
  or published again.

## 0.3.0

- `POST /articles/{id}/reject` sends an article in review back to its author
  with a required note. The author edits it and submits it again; submitting
  clears the note.

## 0.2.0

- Review. Authors submit a draft with `POST /articles/{id}/submit`, and another
  user approves it with `POST /articles/{id}/approve`. Only approved articles
  can be published. Authors cannot review their own articles, and an article
  cannot be edited while it is in review or after it is approved.

## 0.1.0

- Authors create (`POST /articles`), edit (`PUT /articles/{id}`), and publish
  (`POST /articles/{id}/publish`) articles. `GET /articles/{id}` shows any
  article to any user; `GET /published` lists published articles, newest
  first.
- pressroom runs behind the company SSO proxy, which sets `X-Forwarded-User`
  on every request; requests without it are rejected with 401.
- An action the article's current state does not allow, including repeating
  one that already happened, returns 409 and names the current state.
- Articles are kept in memory.
