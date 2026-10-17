package articles

import (
	"cmp"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"
)

// Service holds the articles. Each method runs under one lock, so a check and
// the change it guards cannot interleave with another request.
type Service struct {
	mu       sync.Mutex
	articles map[string]*Article
	nextID   int
	now      func() time.Time
}

func NewService(now func() time.Time) *Service {
	return &Service{articles: map[string]*Article{}, now: now}
}

func (s *Service) Create(author, title, body string) (Article, error) {
	title, body, err := validate(title, body)
	if err != nil {
		return Article{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nextID++
	a := &Article{ID: fmt.Sprintf("a%d", s.nextID), seq: s.nextID, Author: author, Title: title, Body: body, UpdatedAt: s.now()}
	s.articles[a.ID] = a
	return *a, nil
}

func (s *Service) Get(id string) (Article, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	a, ok := s.articles[id]
	if !ok {
		return Article{}, ErrNotFound
	}
	return *a, nil
}

// Published returns the published, unarchived articles, newest first; articles
// published at the same time come in the order they were created.
func (s *Service) Published() []Article {
	s.mu.Lock()
	defer s.mu.Unlock()
	var out []Article
	for _, a := range s.articles {
		if a.Published && !a.Archived {
			out = append(out, *a)
		}
	}
	slices.SortFunc(out, func(a, b Article) int {
		if c := b.PublishedAt.Compare(a.PublishedAt); c != 0 {
			return c
		}
		return cmp.Compare(a.seq, b.seq)
	})
	return out
}

func (s *Service) Edit(id, user, title, body string) (Article, error) {
	title, body, err := validate(title, body)
	if err != nil {
		return Article{}, err
	}
	return s.change(id, user, "edit", author, func(a *Article, now time.Time) error {
		if a.Submitted || a.Approved || a.Published || a.Archived {
			return notAllowed("edit", a)
		}
		a.Title, a.Body = title, body
		return nil
	})
}

func (s *Service) Submit(id, user string) (Article, error) {
	return s.change(id, user, "submit", author, func(a *Article, now time.Time) error {
		if a.Submitted || a.Approved || a.Published || a.Archived {
			return notAllowed("submit", a)
		}
		a.Submitted, a.Rejected = true, false
		a.Reviewer, a.ReviewNote = "", ""
		return nil
	})
}

func (s *Service) Approve(id, user string) (Article, error) {
	return s.change(id, user, "approve", reviewer, func(a *Article, now time.Time) error {
		if !a.Submitted || a.Approved {
			return notAllowed("approve", a)
		}
		a.Approved, a.Reviewer = true, user
		return nil
	})
}

func (s *Service) Reject(id, user, note string) (Article, error) {
	note = strings.TrimSpace(note)
	if note == "" {
		return Article{}, fmt.Errorf("%w: a rejection needs a note", ErrInvalid)
	}
	return s.change(id, user, "reject", reviewer, func(a *Article, now time.Time) error {
		if !a.Submitted || a.Approved {
			return notAllowed("reject", a)
		}
		a.Submitted, a.Rejected = false, true
		a.Reviewer, a.ReviewNote = user, note
		return nil
	})
}

func (s *Service) Publish(id, user string) (Article, error) {
	return s.change(id, user, "publish", author, func(a *Article, now time.Time) error {
		if !a.Approved || a.Published {
			return notAllowed("publish", a)
		}
		a.Published, a.PublishedAt = true, now
		return nil
	})
}

func (s *Service) Archive(id, user string) (Article, error) {
	return s.change(id, user, "archive", author, func(a *Article, now time.Time) error {
		if !a.Published || a.Archived {
			return notAllowed("archive", a)
		}
		a.Archived = true
		return nil
	})
}

type role int

const (
	author role = iota
	reviewer
)

// change checks that user may act on the article in the given role, applies
// fn under the lock with the time of the change, and stamps UpdatedAt with
// that time if fn succeeds.
func (s *Service) change(id, user, action string, r role, fn func(a *Article, now time.Time) error) (Article, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	a, ok := s.articles[id]
	if !ok {
		return Article{}, ErrNotFound
	}
	if r == author && user != a.Author {
		return Article{}, fmt.Errorf("%w: only the author can %s an article", ErrForbidden, action)
	}
	if r == reviewer && user == a.Author {
		return Article{}, fmt.Errorf("%w: authors cannot %s their own articles", ErrForbidden, action)
	}
	now := s.now()
	if err := fn(a, now); err != nil {
		return Article{}, err
	}
	a.UpdatedAt = now
	return *a, nil
}
