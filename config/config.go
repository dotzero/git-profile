package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Entry is the entry in config file
type Entry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Config is the config storage
type Config struct {
	Profiles map[string][]Entry `json:"profiles"`
}

// New initializes and returns a new Config
func New() *Config {
	return &Config{
		Profiles: make(map[string][]Entry),
	}
}

// Len returns number of profiles
func (c *Config) Len() int {
	return len(c.Profiles)
}

// Lookup returns the profile with the given name
func (c *Config) Lookup(name string) ([]Entry, bool) {
	entries, ok := c.Profiles[name]

	return entries, ok
}

// Names returns profile names
func (c *Config) Names() []string {
	names := make([]string, 0, len(c.Profiles))

	for name := range c.Profiles {
		names = append(names, name)
	}

	return names
}

// Delete deletes the value for a key in the profile
func (c *Config) Delete(profile string, value string) bool {
	if _, ok := c.Profiles[profile]; !ok {
		return false
	}

	entries := c.Profiles[profile][:0]

	for _, e := range c.Profiles[profile] {
		if e.Key != value {
			entries = append(entries, e)
		}
	}

	delete(c.Profiles, profile)

	if len(entries) > 0 {
		c.Profiles[profile] = entries
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

	c.Profiles[profile] = append(c.Profiles[profile], Entry{key, value})
}

// Save stores profiles to json file
func (c *Config) Save(filename string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, 0644) // nolint
}

// Load profiles from json file
func (c *Config) Load(filename string) (err error) {
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		err = c.Save(filename)
		if err != nil {
			return err
		}
	}

	body, err := ioutil.ReadFile(filename) // nolint
	if err != nil {
		return err
	}

	return json.Unmarshal(body, c)
}
