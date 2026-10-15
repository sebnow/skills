# Operating bramble

## Running the server

`bramble-server` listens on `-listen` (default `127.0.0.1:8080`) and keeps
everything in the SQLite file given by `-db`. Runners connect to it over
TLS and pull tasks.

## The GUI

The Tasks page lists every task, newest first, with its state and runner.
A task's page shows its log, its artifacts and the actions its state
allows: Stop while it runs, Retry once it has ended, Dismiss once it is
stopped, failed or finished.

The settings page lists the projects and the runners, and can create,
edit and delete projects; a project's own page lists its tasks grouped by
state.

## Projects

A project names the defaults its tasks start with: runner pool, branch
and environment. A task started from a project's page belongs to it.

The project page opens with the overview, then a pane holding
the project's settings, the form that edits the project. Delete sits
under the form and is shown only while no task belongs to the project.

## API

`GET /v1/projects` lists the projects. `POST /v1/projects` creates one
and `PUT /v1/projects/{project}` replaces its settings.
`DELETE /v1/projects/{project}` deletes one (409 while a task belongs
to it).

`GET /v1/tasks` lists tasks; `POST /v1/tasks/{task}/stop` stops one and
`POST /v1/tasks/{task}/dismiss` dismisses one.
