package articles

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	ErrNotFound   = errors.New("article not found")
	ErrInvalid    = errors.New("invalid article")
	ErrForbidden  = errors.New("forbidden")
	ErrNotAllowed = errors.New("not allowed in the article's current state")
)

const maxTitle = 200

type Article struct {
	ID     string
	seq    int // creation order
	Author string
	Title  string
	Body   string

	Submitted bool
	Approved  bool
	Rejected  bool
	Published bool
	Archived  bool

	Reviewer    string
	ReviewNote  string
	UpdatedAt   time.Time
	PublishedAt time.Time
}

// State is the article's lifecycle state as shown by the API.
func (a Article) State() string {
	switch {
	case a.Archived:
		return "archived"
	case a.Published:
		return "published"
	case a.Approved:
		return "approved"
	case a.Rejected:
		return "changes_requested"
	case a.Submitted:
		return "in_review"
	default:
		return "draft"
	}
}

func notAllowed(action string, a *Article) error {
	return fmt.Errorf("%w: cannot %s an article that is %s", ErrNotAllowed, action, a.State())
}

func validate(title, body string) (string, string, error) {
	title, body = strings.TrimSpace(title), strings.TrimSpace(body)
	switch {
	case title == "":
		return "", "", fmt.Errorf("%w: title is required", ErrInvalid)
	case utf8.RuneCountInString(title) > maxTitle:
		return "", "", fmt.Errorf("%w: title is longer than %d characters", ErrInvalid, maxTitle)
	case body == "":
		return "", "", fmt.Errorf("%w: body is required", ErrInvalid)
	}
	return title, body, nil
}
