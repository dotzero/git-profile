package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
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
	log.Println("[DEBUG] GetProfile", profile)
	entries, ok := c.Profiles[profile]
	if ok {
		return entries, true
	}

	return nil, false
}

// RemoveProfile removes a profile
func (c *Config) RemoveProfile(profile string) bool {
	log.Println("[DEBUG] RemoveProfile", profile)
	delete(c.Profiles, profile)
	return true
}

// SetValue adds an entry to profile
func (c *Config) SetValue(profile string, entry Entry) bool {
	log.Println("[DEBUG] SetValue", profile, entry)
	c.RemoveValue(profile, entry.Key)
	c.Profiles[profile] = append(c.Profiles[profile], entry)
	return true
}

// RemoveValue removes an entry from profile
func (c *Config) RemoveValue(profile string, value string) bool {
	log.Println("[DEBUG] RemoveValue", profile, value)
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
func (c *Config) Save(filename string) (err error) {
	log.Println("[DEBUG] Save", filename)

	body, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	err = c.writeFile(filename, body)
	if err != nil {
		return err
	}

	return nil
}

// Load profiles from json file
func (c *Config) Load(filename string) (err error) {
	log.Println("[DEBUG] Load", filename)

	err = c.ensureFile(filename)
	if err != nil {
		return err
	}

	body, err := c.readFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, c)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) ensureFile(filename string) (err error) {
	log.Println("[DEBUG] ensureFile", filename)
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		err = c.Save(filename)
	}
	return
}

func (c *Config) writeFile(filename string, body []byte) (err error) {
	log.Println("[DEBUG] writeFile", filename, string(body))
	err = ioutil.WriteFile(filename, body, 0644)
	return
}

func (c *Config) readFile(filename string) (body []byte, err error) {
	log.Println("[DEBUG] readFile", filename)
	body, err = ioutil.ReadFile(filename)
	return
}
