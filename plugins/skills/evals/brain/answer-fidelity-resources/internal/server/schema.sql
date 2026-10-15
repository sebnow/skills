-- Schema version 1. Later versions are migrations in store.go.

CREATE TABLE projects (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL UNIQUE,
	pool TEXT NOT NULL,
	branch TEXT NOT NULL,
	environment TEXT NOT NULL DEFAULT ''
);

CREATE TABLE runners (
	id TEXT PRIMARY KEY,
	pool TEXT NOT NULL,
	last_seen TEXT
);

CREATE TABLE tasks (
	id TEXT PRIMARY KEY,
	project TEXT REFERENCES projects (id),
	runner TEXT REFERENCES runners (id),
	state TEXT NOT NULL,
	command TEXT NOT NULL,
	created TEXT NOT NULL,
	dismissed TEXT
);

CREATE TABLE events (
	task_id TEXT NOT NULL REFERENCES tasks (id),
	seq INTEGER NOT NULL,
	kind TEXT NOT NULL,
	body TEXT NOT NULL,
	PRIMARY KEY (task_id, seq)
);

CREATE TABLE logs (
	task_id TEXT NOT NULL REFERENCES tasks (id),
	line INTEGER NOT NULL,
	text TEXT NOT NULL,
	PRIMARY KEY (task_id, line)
);
