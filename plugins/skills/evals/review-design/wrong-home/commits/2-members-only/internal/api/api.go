package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/brightloop/deskbook/internal/bookings"
	"github.com/brightloop/deskbook/internal/rooms"
)

const maxBody = 64 << 10

type server struct {
	cfg   rooms.Config
	store *bookings.Store
	now   func() time.Time
}

// Handler serves the API. Every request must carry X-Forwarded-User, set by
// the SSO proxy in front of deskbook.
func Handler(cfg rooms.Config, store *bookings.Store, now func() time.Time) http.Handler {
	s := &server{cfg: cfg, store: store, now: now}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /bookings", s.book)
	mux.HandleFunc("DELETE /bookings/{id}", s.cancel)
	mux.HandleFunc("GET /rooms/{id}/bookings", s.upcoming)
	return requireUser(mux)
}

type bookingRequest struct {
	Room  string    `json:"room"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

func (s *server) book(w http.ResponseWriter, r *http.Request) {
	var req bookingRequest
	if !decode(w, r, &req) {
		return
	}
	switch {
	case !req.End.After(req.Start):
		writeError(w, http.StatusBadRequest, "end must be after start")
		return
	case req.Start.Before(s.now()):
		writeError(w, http.StatusBadRequest, "start is in the past")
		return
	}
	room, ok := s.cfg.Room(req.Room)
	if !ok {
		writeError(w, http.StatusNotFound, "unknown room")
		return
	}
	if room.MembersOnly && !s.cfg.IsMember(user(r)) {
		writeError(w, http.StatusForbidden, fmt.Sprintf("%s is for members only", room.Name))
		return
	}
	if room.MaxMinutes > 0 && req.End.Sub(req.Start) > time.Duration(room.MaxMinutes)*time.Minute {
		writeError(w, http.StatusUnprocessableEntity, fmt.Sprintf("%s can be booked for at most %d minutes", room.Name, room.MaxMinutes))
		return
	}
	b, err := s.store.Add(room.ID, user(r), req.Start, req.End)
	switch {
	case errors.Is(err, bookings.ErrConflict):
		writeError(w, http.StatusConflict, err.Error())
	case err != nil:
		writeError(w, http.StatusInternalServerError, "internal error")
	default:
		writeJSON(w, http.StatusCreated, b)
	}
}

func (s *server) cancel(w http.ResponseWriter, r *http.Request) {
	switch err := s.store.Cancel(r.PathValue("id"), user(r)); {
	case errors.Is(err, bookings.ErrNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, bookings.ErrForbidden):
		writeError(w, http.StatusForbidden, err.Error())
	default:
		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *server) upcoming(w http.ResponseWriter, r *http.Request) {
	room, ok := s.cfg.Room(r.PathValue("id"))
	if !ok {
		writeError(w, http.StatusNotFound, "unknown room")
		return
	}
	list := s.store.Upcoming(room.ID, s.now())
	if list == nil {
		list = []bookings.Booking{}
	}
	writeJSON(w, http.StatusOK, list)
}

func requireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user(r) == "" {
			writeError(w, http.StatusUnauthorized, "missing X-Forwarded-User")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func user(r *http.Request) string { return r.Header.Get("X-Forwarded-User") }

// decode reads exactly one JSON object from the body into v. A body over
// maxBody is refused whatever it contains.
func decode(w http.ResponseWriter, r *http.Request, v any) bool {
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxBody))
	var tooLarge *http.MaxBytesError
	if errors.As(err, &tooLarge) {
		writeError(w, http.StatusRequestEntityTooLarge, "request body too large")
		return false
	}
	if err == nil {
		dec := json.NewDecoder(bytes.NewReader(body))
		dec.DisallowUnknownFields()
		if err = dec.Decode(v); err == nil && dec.Decode(&struct{}{}) != io.EOF {
			err = errors.New("body must hold a single JSON object")
		}
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body: "+err.Error())
		return false
	}
	return true
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
