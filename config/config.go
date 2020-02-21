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
	return &Config{}
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

// Store sets the value for a key in the profile
func (c *Config) Store(profile string, key string, value string) {
	c.Delete(profile, key)

	if c.Profiles == nil {
		c.Profiles = make(map[string][]Entry)
	}

	c.Profiles[profile] = append(c.Profiles[profile], Entry{key, value})

}

// Save stores profiles to json file
func (c *Config) Save(filename string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, 0644)
}

// Load profiles from json file
func (c *Config) Load(filename string) (err error) {
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		err = c.Save(filename)
		if err != nil {
			return err
		}
	}

	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, c)
}
