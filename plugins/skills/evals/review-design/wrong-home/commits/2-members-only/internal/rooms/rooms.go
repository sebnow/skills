package rooms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

type Room struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// MaxMinutes caps how long one booking of the room may run; 0 means no cap.
	MaxMinutes int `json:"max_minutes"`
	// MembersOnly rooms can be booked only by users listed in Config.Members.
	MembersOnly bool `json:"members_only"`
}

type Config struct {
	Rooms   []Room   `json:"rooms"`
	Members []string `json:"members"`
}

// maxMinutesLimit is the largest cap a room may set: one week.
const maxMinutesLimit = 7 * 24 * 60

// Load reads the room configuration from a JSON file.
func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	var c Config
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&c); err != nil {
		return Config{}, fmt.Errorf("%s: %w", path, err)
	}
	seen := map[string]bool{}
	for _, r := range c.Rooms {
		switch {
		case r.ID == "" || r.Name == "":
			return Config{}, fmt.Errorf("%s: every room needs an id and a name", path)
		case seen[r.ID]:
			return Config{}, fmt.Errorf("%s: room %q is listed twice", path, r.ID)
		case r.MembersOnly && len(c.Members) == 0:
			return Config{}, fmt.Errorf("%s: room %q is members_only but no members are listed", path, r.ID)
		case r.MaxMinutes < 0 || r.MaxMinutes > maxMinutesLimit:
			return Config{}, fmt.Errorf("%s: room %q: max_minutes must be between 0 and %d", path, r.ID, maxMinutesLimit)
		}
		seen[r.ID] = true
	}
	members := map[string]bool{}
	for _, m := range c.Members {
		switch {
		case m == "":
			return Config{}, fmt.Errorf("%s: members must not be empty", path)
		case members[m]:
			return Config{}, fmt.Errorf("%s: member %q is listed twice", path, m)
		}
		members[m] = true
	}
	return c, nil
}

func (c Config) Room(id string) (Room, bool) {
	for _, r := range c.Rooms {
		if r.ID == id {
			return r, true
		}
	}
	return Room{}, false
}

func (c Config) IsMember(user string) bool {
	for _, m := range c.Members {
		if m == user {
			return true
		}
	}
	return false
}
