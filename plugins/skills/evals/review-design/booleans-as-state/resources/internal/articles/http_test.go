package articles

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestStatuses(t *testing.T) {
	h := Handler(NewService(time.Now))
	call := func(method, path, user, body string) int {
		req := httptest.NewRequest(method, path, strings.NewReader(body))
		if user != "" {
			req.Header.Set("X-Forwarded-User", user)
		}
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		return rec.Code
	}
	for _, tc := range []struct {
		method, path, user, body string
		want                     int
	}{
		{"GET", "/published", "", "", http.StatusUnauthorized},
		{"POST", "/articles", "ana", `{"title": "Menu", "body": "Soup."}`, http.StatusCreated},
		{"POST", "/articles", "ana", `{"title": "", "body": "Soup."}`, http.StatusBadRequest},
		{"POST", "/articles", "ana", `{"title": "Menu", "body": "Soup.", "tags": []}`, http.StatusBadRequest},
		{"POST", "/articles/a1/publish", "ana", "", http.StatusConflict},
		{"POST", "/articles/a1/submit", "ben", "", http.StatusForbidden},
		{"POST", "/articles/a1/submit", "ana", "", http.StatusOK},
		{"POST", "/articles/a1/reject", "ben", `{"note": ""}`, http.StatusBadRequest},
		{"POST", "/articles/a1/approve", "ben", "", http.StatusOK},
		{"GET", "/articles/nope", "ana", "", http.StatusNotFound},
		{"POST", "/articles", "ana", `{"title": "Menu", "body": "Soup."} {}`, http.StatusBadRequest},
		{"POST", "/articles", "ana", `{"title": "Menu", "body": "` + strings.Repeat("x", maxBody) + `"}`, http.StatusRequestEntityTooLarge},
		{"POST", "/articles", "ana", `{"title": "Menu", "body": "Soup."}` + strings.Repeat(" ", maxBody), http.StatusRequestEntityTooLarge},
		{"POST", "/articles", "ana", `{"title": "Menu", "body": "Soup."}x` + strings.Repeat("x", maxBody), http.StatusRequestEntityTooLarge},
	} {
		if got := call(tc.method, tc.path, tc.user, tc.body); got != tc.want {
			t.Errorf("%s %s as %q: status %d, want %d", tc.method, tc.path, tc.user, got, tc.want)
		}
	}
}

func TestErrorsAreJSON(t *testing.T) {
	h := Handler(NewService(time.Now))
	req := httptest.NewRequest("GET", "/articles/nope", nil)
	req.Header.Set("X-Forwarded-User", "ana")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	var body struct{ Error string }
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil || body.Error == "" {
		t.Fatalf("body %q: want a JSON object with an error message (%v)", rec.Body.String(), err)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", ct)
	}
}
