package config

import (
	"encoding/json"
	"io/ioutil"
)

// Entry is the entry in config file
type Entry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Config is the config file
type Config struct {
	Profiles map[string][]Entry `json:"profiles"`
}

// NewConfig initializes and returns a new Config
func NewConfig() *Config {
	var c Config
	c.Profiles = make(map[string][]Entry)
	return &c
}

// GetProfile gets a profile entries
func (c *Config) GetProfile(profile string) ([]Entry, bool) {
	entries, ok := c.Profiles[profile]
	if ok {
		return entries, true
	}

	return nil, false
}

// RemoveProfile removes a profile
func (c *Config) RemoveProfile(profile string) bool {
	delete(c.Profiles, profile)
	return true
}

// SetValue adds an entry to profile
func (c *Config) SetValue(profile string, entry Entry) bool {
	c.RemoveValue(profile, entry.Key)
	c.Profiles[profile] = append(c.Profiles[profile], entry)
	return true
}

// RemoveValue removes an entry from profile
func (c *Config) RemoveValue(profile string, value string) bool {
	if _, ok := c.Profiles[profile]; !ok {
		return false
	}

	entries := c.Profiles[profile][:0]
	for _, e := range c.Profiles[profile] {
		if e.Key != value {
			entries = append(entries, e)
		}
	}

	c.RemoveProfile(profile)
	if len(entries) > 0 {
		c.Profiles[profile] = entries
	}

	return true
}

// Save profiles to json file
func (c *Config) Save(filename string) (bool, error) {
	var err error

	str, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return false, err
	}

	err = ioutil.WriteFile(filename, str, 0644)
	if err != nil {
		return false, err
	}

	return true, nil
}

// Load profiles from json file
func (c *Config) Load(filename string) (bool, error) {
	var err error

	str, err := ioutil.ReadFile(filename)
	if err != nil {
		return false, err
	}

	err = json.Unmarshal(str, c)
	if err != nil {
		return false, err
	}

	return true, nil
}
