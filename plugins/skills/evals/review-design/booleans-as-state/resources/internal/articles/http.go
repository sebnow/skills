package articles

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

const maxBody = 1 << 20

type articleView struct {
	ID          string     `json:"id"`
	Author      string     `json:"author"`
	Title       string     `json:"title"`
	Body        string     `json:"body"`
	State       string     `json:"state"`
	Reviewer    string     `json:"reviewer,omitempty"`
	ReviewNote  string     `json:"review_note,omitempty"`
	UpdatedAt   time.Time  `json:"updated_at"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
}

func view(a Article) articleView {
	v := articleView{ID: a.ID, Author: a.Author, Title: a.Title, Body: a.Body, State: a.State(), UpdatedAt: a.UpdatedAt}
	if a.Approved || a.Rejected {
		v.Reviewer = a.Reviewer
	}
	if a.Rejected {
		v.ReviewNote = a.ReviewNote
	}
	if a.Published {
		t := a.PublishedAt
		v.PublishedAt = &t
	}
	return v
}

type content struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// Handler serves the API. Every request must carry X-Forwarded-User, set by
// the SSO proxy in front of pressroom.
func Handler(s *Service) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /articles", func(w http.ResponseWriter, r *http.Request) {
		var c content
		if !decode(w, r, &c) {
			return
		}
		a, err := s.Create(user(r), c.Title, c.Body)
		respond(w, http.StatusCreated, a, err)
	})
	mux.HandleFunc("GET /articles/{id}", func(w http.ResponseWriter, r *http.Request) {
		a, err := s.Get(r.PathValue("id"))
		respond(w, http.StatusOK, a, err)
	})
	mux.HandleFunc("PUT /articles/{id}", func(w http.ResponseWriter, r *http.Request) {
		var c content
		if !decode(w, r, &c) {
			return
		}
		a, err := s.Edit(r.PathValue("id"), user(r), c.Title, c.Body)
		respond(w, http.StatusOK, a, err)
	})
	mux.HandleFunc("POST /articles/{id}/reject", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Note string `json:"note"`
		}
		if !decode(w, r, &req) {
			return
		}
		a, err := s.Reject(r.PathValue("id"), user(r), req.Note)
		respond(w, http.StatusOK, a, err)
	})
	for action, fn := range map[string]func(id, user string) (Article, error){
		"submit":  s.Submit,
		"approve": s.Approve,
		"publish": s.Publish,
		"archive": s.Archive,
	} {
		mux.HandleFunc("POST /articles/{id}/"+action, func(w http.ResponseWriter, r *http.Request) {
			a, err := fn(r.PathValue("id"), user(r))
			respond(w, http.StatusOK, a, err)
		})
	}
	mux.HandleFunc("GET /published", func(w http.ResponseWriter, r *http.Request) {
		views := []articleView{}
		for _, a := range s.Published() {
			views = append(views, view(a))
		}
		writeJSON(w, http.StatusOK, views)
	})
	return requireUser(mux)
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

func respond(w http.ResponseWriter, status int, a Article, err error) {
	switch {
	case err == nil:
		writeJSON(w, status, view(a))
	case errors.Is(err, ErrInvalid):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, ErrForbidden):
		writeError(w, http.StatusForbidden, err.Error())
	case errors.Is(err, ErrNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, ErrNotAllowed):
		writeError(w, http.StatusConflict, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "internal error")
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
