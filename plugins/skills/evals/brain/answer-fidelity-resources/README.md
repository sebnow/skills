# bramble

A self-hosted server that runs build and test tasks on runner machines.
Tasks belong to projects; a project holds the defaults its tasks start
with (runner pool, branch, environment).

```
bramble-server -db bramble.db -listen 127.0.0.1:8080
```

Runners are a separate program that connects to the server and pulls
tasks; this repository holds the server.

See `docs/operating.md` for the GUI and API, and `docs/adr/` for the
records behind them.
