---
name: reversible-decision
tags: [adr-writing, reversible]
max_turns: 12
allowed_tools: [Skill, Agent]
runs: 3
---

Can you write an ADR for switching the dev log format from JSON to plain text? In `config/dev.env` we have `LOG_FORMAT=json`, and tailing `docker compose logs -f notifier worker` is unreadable: every line is a 400-character JSON blob with `ts`, `level`, `msg`, `trace_id`, `span_id`, and whatever fields the handler added. The logger (zerolog) already supports `LOG_FORMAT=text`, which gives `12:04:31 INF sent reminder user=8812 channel=email`. So the change is just that one line in `config/dev.env`. Staging and prod stay on JSON because Loki parses those. Put the ADR in your reply.
