// Package vessel reads the vessel file: what tidewright knows about the
// vessel that no log carries.
package vessel

import "github.com/BurntSushi/toml"

// File is vessel.toml.
type File struct {
	Name       string   `toml:"name"`
	HomePort   string   `toml:"home_port"`
	Regulation struct {
		Amendments string `toml:"amendments"`
	} `toml:"regulation"`
	Class struct {
		Society   string `toml:"society"`
		ReportURL string `toml:"report_url"`
	} `toml:"class"`
	Harbour struct {
		Register string `toml:"register"`
		Berth    string `toml:"berth"`
	} `toml:"harbour"`
	Engines []Engine `toml:"engine"`
}

// Engine is one [[engine]] entry.
type Engine struct {
	ID        string `toml:"id"`
	Model     string `toml:"model"`
	Telemetry string `toml:"telemetry"`
}

// Load reads the vessel file at path.
func Load(path string) (File, error) {
	var f File
	_, err := toml.DecodeFile(path, &f)
	return f, err
}
