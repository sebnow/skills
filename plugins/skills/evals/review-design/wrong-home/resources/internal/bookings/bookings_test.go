package bookings

import (
	"errors"
	"testing"
	"time"
)

var t0 = time.Date(2026, 3, 2, 9, 0, 0, 0, time.UTC)

func TestOverlappingBookingIsRefused(t *testing.T) {
	s := NewStore()
	if _, err := s.Add("garden", "ana", t0, t0.Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Add("garden", "ben", t0.Add(30*time.Minute), t0.Add(90*time.Minute)); !errors.Is(err, ErrConflict) {
		t.Fatalf("err = %v, want ErrConflict", err)
	}
	if _, err := s.Add("garden", "ben", t0.Add(time.Hour), t0.Add(2*time.Hour)); err != nil {
		t.Fatalf("back-to-back booking: %v", err)
	}
}

func TestOnlyTheMemberCanCancel(t *testing.T) {
	s := NewStore()
	b, _ := s.Add("garden", "ana", t0, t0.Add(time.Hour))
	if err := s.Cancel(b.ID, "ben"); !errors.Is(err, ErrForbidden) {
		t.Fatalf("err = %v, want ErrForbidden", err)
	}
	if err := s.Cancel(b.ID, "ana"); err != nil {
		t.Fatal(err)
	}
	if got := s.Upcoming("garden", t0); len(got) != 0 {
		t.Fatalf("upcoming = %v after cancel", got)
	}
}

func TestUpcomingSkipsEndedBookingsAndSortsByStart(t *testing.T) {
	s := NewStore()
	late, _ := s.Add("garden", "ana", t0.Add(3*time.Hour), t0.Add(4*time.Hour))
	s.Add("garden", "ana", t0, t0.Add(time.Hour))
	early, _ := s.Add("garden", "ben", t0.Add(time.Hour), t0.Add(2*time.Hour))

	got := s.Upcoming("garden", t0.Add(time.Hour))
	if len(got) != 2 || got[0].ID != early.ID || got[1].ID != late.ID {
		t.Fatalf("upcoming = %v", got)
	}
}
