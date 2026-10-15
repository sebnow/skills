package server

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"

	_ "modernc.org/sqlite"
)

// handler serves the GUI and the API from one Store.
type handler struct {
	store *Store
	log   *slog.Logger
}

// fail logs an unexpected error and answers 500 without its text.
func (h *handler) fail(w http.ResponseWriter, err error) {
	h.log.Error("request failed", "error", err)
	http.Error(w, "Something went wrong.", http.StatusInternalServerError)
}

// Run opens the database at dbPath, migrates it and serves on listen
// until the listener fails.
func Run(dbPath, listen string) error {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("open %s: %w", dbPath, err)
	}
	s := &Store{db: db}
	if err := s.migrate(context.Background()); err != nil {
		return err
	}
	h := &handler{store: s, log: slog.Default()}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /projects/{project}/delete", h.postDeleteProject)
	return http.ListenAndServe(listen, mux)
}
