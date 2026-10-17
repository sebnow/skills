package bookings

import (
	"cmp"
	"errors"
	"fmt"
	"slices"
	"sync"
	"time"
)

var (
	ErrNotFound  = errors.New("booking not found")
	ErrForbidden = errors.New("only the member who made a booking can cancel it")
	ErrConflict  = errors.New("the room is already booked for part of that time")
)

type Booking struct {
	ID     string    `json:"id"`
	Room   string    `json:"room"`
	Member string    `json:"member"`
	Start  time.Time `json:"start"`
	End    time.Time `json:"end"`
	seq    int
}

// Store holds the bookings. Each method runs under one lock, so a check and
// the change it guards cannot interleave with another request.
type Store struct {
	mu       sync.Mutex
	bookings map[string]*Booking
	next     int
}

func NewStore() *Store { return &Store{bookings: map[string]*Booking{}} }

// Add books the room from start to end for member, unless another booking of
// the room overlaps that time.
func (s *Store) Add(room, member string, start, end time.Time) (Booking, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.conflicts(room, start, end) {
		return Booking{}, ErrConflict
	}
	s.next++
	b := &Booking{ID: fmt.Sprintf("b%d", s.next), Room: room, Member: member, Start: start.UTC(), End: end.UTC(), seq: s.next}
	s.bookings[b.ID] = b
	return *b, nil
}

func (s *Store) conflicts(room string, start, end time.Time) bool {
	for _, b := range s.bookings {
		if b.Room == room && b.Start.Before(end) && start.Before(b.End) {
			return true
		}
	}
	return false
}

// Upcoming returns the room's bookings that have not ended, earliest first.
func (s *Store) Upcoming(room string, now time.Time) []Booking {
	s.mu.Lock()
	defer s.mu.Unlock()
	var out []Booking
	for _, b := range s.bookings {
		if b.Room == room && b.End.After(now) {
			out = append(out, *b)
		}
	}
	slices.SortFunc(out, func(a, b Booking) int {
		if c := a.Start.Compare(b.Start); c != 0 {
			return c
		}
		return cmp.Compare(a.seq, b.seq)
	})
	return out
}

func (s *Store) Cancel(id, member string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	b, ok := s.bookings[id]
	switch {
	case !ok:
		return ErrNotFound
	case b.Member != member:
		return ErrForbidden
	}
	delete(s.bookings, id)
	return nil
}
