package articles

import (
	"errors"
	"testing"
	"time"
)

type clock struct{ t time.Time }

func (c *clock) now() time.Time { c.t = c.t.Add(time.Minute); return c.t }

func newService() *Service {
	c := &clock{t: time.Date(2026, 4, 1, 9, 0, 0, 0, time.UTC)}
	return NewService(c.now)
}

func draft(t *testing.T, s *Service) Article {
	t.Helper()
	a, err := s.Create("ana", "Spring menu", "Asparagus is back.")
	if err != nil {
		t.Fatal(err)
	}
	return a
}

// ok returns a function that fails the test on error, for wrapping calls
// such as ok(t)(s.Submit(id, user)).
func ok(t *testing.T) func(Article, error) Article {
	return func(a Article, err error) Article {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
		return a
	}
}

func TestArticleGoesFromDraftToArchived(t *testing.T) {
	s := newService()
	id := draft(t, s).ID
	steps := []struct {
		do   func() (Article, error)
		want string
	}{
		{func() (Article, error) { return s.Submit(id, "ana") }, "in_review"},
		{func() (Article, error) { return s.Approve(id, "ben") }, "approved"},
		{func() (Article, error) { return s.Publish(id, "ana") }, "published"},
		{func() (Article, error) { return s.Archive(id, "ana") }, "archived"},
	}
	for _, step := range steps {
		if got := ok(t)(step.do()).State(); got != step.want {
			t.Fatalf("state = %q, want %q", got, step.want)
		}
	}
}

func TestRejectedArticleIsEditedAndResubmitted(t *testing.T) {
	s := newService()
	id := draft(t, s).ID
	ok(t)(s.Submit(id, "ana"))
	a := ok(t)(s.Reject(id, "ben", "Needs prices."))
	if a.State() != "changes_requested" || a.ReviewNote != "Needs prices." {
		t.Fatalf("got %q with note %q", a.State(), a.ReviewNote)
	}
	ok(t)(s.Edit(id, "ana", "Spring menu", "Asparagus is back, 14 EUR."))
	a = ok(t)(s.Submit(id, "ana"))
	if a.State() != "in_review" || a.ReviewNote != "" || a.Reviewer != "" {
		t.Fatalf("got %q with reviewer %q and note %q", a.State(), a.Reviewer, a.ReviewNote)
	}
}

func TestActionsOutOfOrderAreRefused(t *testing.T) {
	for _, tc := range []struct {
		name  string
		setup []string
		do    func(s *Service, id string) (Article, error)
	}{
		{"publish a draft", nil, func(s *Service, id string) (Article, error) { return s.Publish(id, "ana") }},
		{"approve a draft", nil, func(s *Service, id string) (Article, error) { return s.Approve(id, "ben") }},
		{"archive a draft", nil, func(s *Service, id string) (Article, error) { return s.Archive(id, "ana") }},
		{"submit twice", []string{"submit"}, func(s *Service, id string) (Article, error) { return s.Submit(id, "ana") }},
		{"edit in review", []string{"submit"}, func(s *Service, id string) (Article, error) { return s.Edit(id, "ana", "t", "b") }},
		{"publish twice", []string{"submit", "approve", "publish"}, func(s *Service, id string) (Article, error) { return s.Publish(id, "ana") }},
		{"submit after archive", []string{"submit", "approve", "publish", "archive"}, func(s *Service, id string) (Article, error) { return s.Submit(id, "ana") }},
	} {
		s := newService()
		id := draft(t, s).ID
		steps := map[string]func() (Article, error){
			"submit":  func() (Article, error) { return s.Submit(id, "ana") },
			"approve": func() (Article, error) { return s.Approve(id, "ben") },
			"publish": func() (Article, error) { return s.Publish(id, "ana") },
			"archive": func() (Article, error) { return s.Archive(id, "ana") },
		}
		for _, step := range tc.setup {
			ok(t)(steps[step]())
		}
		if _, err := tc.do(s, id); !errors.Is(err, ErrNotAllowed) {
			t.Errorf("%s: err = %v, want ErrNotAllowed", tc.name, err)
		}
	}
}

func TestAuthorCannotReviewOwnArticle(t *testing.T) {
	s := newService()
	id := draft(t, s).ID
	ok(t)(s.Submit(id, "ana"))
	if _, err := s.Approve(id, "ana"); !errors.Is(err, ErrForbidden) {
		t.Fatalf("err = %v, want ErrForbidden", err)
	}
}

func TestPublishedListIsNewestFirstWithoutArchived(t *testing.T) {
	s := newService()
	var ids []string
	for range 3 {
		id := draft(t, s).ID
		ok(t)(s.Submit(id, "ana"))
		ok(t)(s.Approve(id, "ben"))
		ok(t)(s.Publish(id, "ana"))
		ids = append(ids, id)
	}
	ok(t)(s.Archive(ids[1], "ana"))

	got := s.Published()
	if len(got) != 2 || got[0].ID != ids[2] || got[1].ID != ids[0] {
		t.Fatalf("published = %v", got)
	}
}

func TestPublishStampsOneTime(t *testing.T) {
	s := newService()
	id := draft(t, s).ID
	ok(t)(s.Submit(id, "ana"))
	ok(t)(s.Approve(id, "ben"))
	a := ok(t)(s.Publish(id, "ana"))
	if !a.PublishedAt.Equal(a.UpdatedAt) {
		t.Fatalf("published_at %v, updated_at %v; want the same time", a.PublishedAt, a.UpdatedAt)
	}
}

func TestArticlesPublishedTogetherKeepCreationOrder(t *testing.T) {
	at := time.Date(2026, 4, 1, 9, 0, 0, 0, time.UTC)
	s := NewService(func() time.Time { return at })
	var ids []string
	for range 12 {
		id := draft(t, s).ID
		ok(t)(s.Submit(id, "ana"))
		ok(t)(s.Approve(id, "ben"))
		ok(t)(s.Publish(id, "ana"))
		ids = append(ids, id)
	}
	for i, a := range s.Published() {
		if a.ID != ids[i] {
			t.Fatalf("position %d: %s, want %s", i, a.ID, ids[i])
		}
	}
}
