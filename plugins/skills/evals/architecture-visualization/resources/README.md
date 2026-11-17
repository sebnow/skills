# inkwell

Converts customer documents (DOCX, ODT, Markdown, images) to PDF and
thumbnails, for the Inkwell public API.

## Running locally

    docker compose -f deploy/compose.yaml up

The API listens on :8080. Create a key with
`python -m inkwell.app create-key --tenant acme`.

## Processes

    python -m inkwell.app serve     # HTTP API
    python -m inkwell.app worker    # conversion workers
    python -m inkwell.app janitor   # nightly cleanup (run from cron)

## Tests

    pytest
