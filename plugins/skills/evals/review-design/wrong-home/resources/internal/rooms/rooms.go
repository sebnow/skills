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
}

type Config struct {
	Rooms []Room `json:"rooms"`
}

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
		}
		seen[r.ID] = true
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
