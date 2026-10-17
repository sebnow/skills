package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/brightloop/deskbook/internal/bookings"
	"github.com/brightloop/deskbook/internal/rooms"
)

var now = time.Date(2026, 3, 2, 8, 0, 0, 0, time.UTC)

func newHandler() http.Handler {
	cfg := rooms.Config{
		Rooms: []rooms.Room{
			{ID: "garden", Name: "Garden room"},
			{ID: "booth-1", Name: "Phone booth 1", MaxMinutes: 60},
			{ID: "boardroom", Name: "Boardroom", MembersOnly: true},
		},
		Members: []string{"ana"},
	}
	return Handler(cfg, bookings.NewStore(), func() time.Time { return now })
}

func call(h http.Handler, method, path, user, body string) int {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if user != "" {
		req.Header.Set("X-Forwarded-User", user)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec.Code
}

func booking(room, start, end string) string {
	return `{"room": "` + room + `", "start": "` + start + `", "end": "` + end + `"}`
}

func TestBookingStatuses(t *testing.T) {
	h := newHandler()
	for _, tc := range []struct {
		name, method, path, user, body string
		want                           int
	}{
		{"no user", "POST", "/bookings", "", booking("garden", "2026-03-02T09:00:00Z", "2026-03-02T10:00:00Z"), http.StatusUnauthorized},
		{"booked", "POST", "/bookings", "ana", booking("garden", "2026-03-02T09:00:00Z", "2026-03-02T10:00:00Z"), http.StatusCreated},
		{"overlap", "POST", "/bookings", "ben", booking("garden", "2026-03-02T09:30:00Z", "2026-03-02T10:30:00Z"), http.StatusConflict},
		{"end before start", "POST", "/bookings", "ana", booking("garden", "2026-03-02T11:00:00Z", "2026-03-02T10:00:00Z"), http.StatusBadRequest},
		{"in the past", "POST", "/bookings", "ana", booking("garden", "2026-03-01T09:00:00Z", "2026-03-01T10:00:00Z"), http.StatusBadRequest},
		{"unknown room", "POST", "/bookings", "ana", booking("attic", "2026-03-02T09:00:00Z", "2026-03-02T10:00:00Z"), http.StatusNotFound},
		{"booth for an hour", "POST", "/bookings", "ana", booking("booth-1", "2026-03-02T09:00:00Z", "2026-03-02T10:00:00Z"), http.StatusCreated},
		{"booth for longer", "POST", "/bookings", "ana", booking("booth-1", "2026-03-02T11:00:00Z", "2026-03-02T12:01:00Z"), http.StatusUnprocessableEntity},
		{"boardroom by a member", "POST", "/bookings", "ana", booking("boardroom", "2026-03-02T09:00:00Z", "2026-03-02T10:00:00Z"), http.StatusCreated},
		{"boardroom by a guest", "POST", "/bookings", "ben", booking("boardroom", "2026-03-02T11:00:00Z", "2026-03-02T12:00:00Z"), http.StatusForbidden},
		{"unknown field", "POST", "/bookings", "ana", `{"room": "garden", "seats": 4}`, http.StatusBadRequest},
		{"cancel by another member", "DELETE", "/bookings/b1", "ben", "", http.StatusForbidden},
		{"cancel", "DELETE", "/bookings/b1", "ana", "", http.StatusNoContent},
	} {
		if got := call(h, tc.method, tc.path, tc.user, tc.body); got != tc.want {
			t.Errorf("%s: status %d, want %d", tc.name, got, tc.want)
		}
	}
}

func TestCheckAnswersLikeBookingWithoutBooking(t *testing.T) {
	h := newHandler()
	call(h, "POST", "/bookings", "ana", booking("garden", "2026-03-02T09:00:00Z", "2026-03-02T10:00:00Z"))
	for _, tc := range []struct {
		name, user, body string
		want             int
	}{
		{"free", "ben", booking("garden", "2026-03-02T10:00:00Z", "2026-03-02T11:00:00Z"), http.StatusOK},
		{"overlap", "ben", booking("garden", "2026-03-02T09:30:00Z", "2026-03-02T10:30:00Z"), http.StatusConflict},
		{"booth too long", "ana", booking("booth-1", "2026-03-02T11:00:00Z", "2026-03-02T12:01:00Z"), http.StatusUnprocessableEntity},
		{"boardroom by a guest", "ben", booking("boardroom", "2026-03-02T11:00:00Z", "2026-03-02T12:00:00Z"), http.StatusForbidden},
		{"unknown room", "ana", booking("attic", "2026-03-02T09:00:00Z", "2026-03-02T10:00:00Z"), http.StatusNotFound},
	} {
		if got := call(h, "POST", "/bookings/check", tc.user, tc.body); got != tc.want {
			t.Errorf("%s: status %d, want %d", tc.name, got, tc.want)
		}
	}
	if got := call(h, "POST", "/bookings", "ben", booking("garden", "2026-03-02T10:00:00Z", "2026-03-02T11:00:00Z")); got != http.StatusCreated {
		t.Fatalf("booking the slot check said was free: status %d", got)
	}
}
