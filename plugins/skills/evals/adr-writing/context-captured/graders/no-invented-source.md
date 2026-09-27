---
type: llm
focus: last_message
weight: 1
---
The conversation before this ADR chose a pg_trgm GIN index over OpenSearch for catalog-api's `/v2/suggest` endpoint. The only files, tickets and documents it names are the load-test script `loadtest/suggest.js`, the ticket CAT-412, and a Slack thread with the platform team. A reader following a reference the ADR cites must find something that exists; any other path, ticket or document the ADR cites was made up by the writer.

Judge only the ADR in the reply. Names of tables, columns, indexes, endpoints, services and tools (`products`, `lower(name)`, `/v2/suggest`, catalog-api, k6, pg_trgm) are not citations and do not count.

PASS if every file path, URL, ticket or issue id, ADR number, dashboard, runbook, report or other document the ADR cites or links to appears in the conversation (`loadtest/suggest.js`, CAT-412).

FAIL if the ADR cites any path, ticket, or document that does not appear in the conversation, for example a migration file path, a results or report file, a dashboard or runbook link, another ADR by number, a URL to external documentation, or a ticket id other than CAT-412.
