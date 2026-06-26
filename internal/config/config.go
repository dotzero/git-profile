package config

import (
	"encoding/json"
	"os"
	"sort"

	"github.com/dotzero/git-profile/internal/oldconfig"
)

// Entry is the entry in config file
type Entry map[string]string

// Config is the config storage
type Config struct {
	Profiles map[string]Entry `json:"profiles"`
}

// New initializes and returns a new Config
func New() *Config {
	return &Config{
		Profiles: make(map[string]Entry),
	}
}

// Len returns number of profiles
func (c *Config) Len() int {
	return len(c.Profiles)
}

// Lookup returns the profile with the given name
func (c *Config) Lookup(name string) (Entry, bool) {
	entries, ok := c.Profiles[name]

	return entries, ok
}

// Names returns profile names
func (c *Config) Names() []string {
	names := make([]string, 0, len(c.Profiles))

	for name := range c.Profiles {
		names = append(names, name)
	}

	sort.Strings(names)

	return names
}

// Delete deletes the value for a key in the profile
func (c *Config) Delete(profile string, key string) bool {
	if _, ok := c.Profiles[profile]; !ok {
		return false
	}

	delete(c.Profiles[profile], key)

	if len(c.Profiles[profile]) == 0 {
		delete(c.Profiles, profile)
	}

	return true
}

// DeleteProfile deletes the profile
func (c *Config) DeleteProfile(profile string) bool {
	if _, ok := c.Profiles[profile]; !ok {
		return false
	}

	delete(c.Profiles, profile)

	return true
}

// Store sets the value for a key in the profile
func (c *Config) Store(profile string, key string, value string) {
	c.Delete(profile, key)

	if _, ok := c.Profiles[profile]; !ok {
		c.Profiles[profile] = make(Entry)
	}

	c.Profiles[profile][key] = value
}

// Save stores profiles to json file
func (c *Config) Save(filename string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0o644) //nolint:gosec
}

// Load profiles from json file
func (c *Config) Load(filename string) (err error) {
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

	// Try to unmarshal as new format (map)
	err = json.Unmarshal(body, c)
	if err == nil {
		return nil // Success with new format
	}

	// Try to unmarshal as old format (array)
	oldCfg := oldconfig.New()
	err = json.Unmarshal(body, oldCfg)
	if err != nil {
		// Neither format worked, return original error
		return json.Unmarshal(body, c)
	}

	// Convert old format to new format
	for profileName, oldEntries := range oldCfg.Profiles {
		for _, entry := range oldEntries {
			c.Store(profileName, entry.Key, entry.Value)
		}
	}

	// Auto-save converted format
	return c.Save(filename)
}
