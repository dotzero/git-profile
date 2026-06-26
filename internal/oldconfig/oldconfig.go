package oldconfig

import (
	"encoding/json"
	"os"
)

// Entry is the entry in config file
type OldEntry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Config is the config storage
type OldConfig struct {
	Profiles map[string][]OldEntry `json:"profiles"`
}

// New initializes and returns a new Config
func New() *OldConfig {
	return &OldConfig{
		Profiles: make(map[string][]OldEntry),
	}
}

// Len returns number of profiles
func (c *OldConfig) Len() int {
	return len(c.Profiles)
}

// Save stores profiles to json file
func (c *OldConfig) Save(filename string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0o644) //nolint:gosec
}

// Load profiles from json file
func (c *OldConfig) Load(filename string) (err error) {
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		err = c.Save(filename)
		if err != nil {
			return err
		}
	}

	body, err := os.ReadFile(filename) //nolint:gosec
	if err != nil {
		return err
	}

	return json.Unmarshal(body, c)
}
