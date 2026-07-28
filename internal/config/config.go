package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/dotzero/git-profile/internal/oldconfig"
)

const (
	xdgAppDir      = "git-profile"
	xdgConfigFile  = "config.json"
	legacyFileName = ".gitprofile"
)

// Entry is the entry in config file
type Entry map[string]string

// Config is the config storage
type Config struct {
	Profiles map[string]Entry `json:"profiles"`

	// legacy is true when the loaded file uses the array-based format.
	// Legacy files keep that format on save.
	legacy bool
}

// New initializes and returns a new Config
func New() *Config {
	return &Config{
		Profiles: make(map[string]Entry),
	}
}

// DefaultPath resolves the default config file path.
// It prefers $XDG_CONFIG_HOME/git-profile/config.json when that file exists.
// Otherwise it falls back to ~/.gitprofile when that file exists.
// If neither file exists, it returns the XDG path for new installs.
func DefaultPath() (string, error) {
	xdgPath, err := xdgConfigPath()
	if err != nil {
		return "", err
	}

	if fileExists(xdgPath) {
		return xdgPath, nil
	}

	legacyPath, err := legacyConfigPath()
	if err != nil {
		return "", err
	}

	if fileExists(legacyPath) {
		return legacyPath, nil
	}

	return xdgPath, nil
}

// XDGPath returns $XDG_CONFIG_HOME/git-profile/config.json.
func XDGPath() (string, error) {
	return xdgConfigPath()
}

// LegacyPath returns ~/.gitprofile.
func LegacyPath() (string, error) {
	return legacyConfigPath()
}

// Migrate copies ~/.gitprofile to the XDG path in map format.
// It does not remove the legacy file. If the XDG file already exists,
// set force to overwrite it.
func Migrate(force bool) (src string, dst string, err error) {
	src, err = legacyConfigPath()
	if err != nil {
		return "", "", err
	}

	dst, err = xdgConfigPath()
	if err != nil {
		return "", "", err
	}

	if !fileExists(src) {
		return "", "", fmt.Errorf("legacy config not found: %s", src)
	}

	if fileExists(dst) && !force {
		return "", "", fmt.Errorf("XDG config already exists: %s (use --force to overwrite)", dst)
	}

	cfg := New()

	err = cfg.Load(src)
	if err != nil {
		return "", "", err
	}

	// Always write the map format to the XDG path.
	cfg.legacy = false

	err = cfg.Save(dst)
	if err != nil {
		return "", "", err
	}

	return src, dst, nil
}

// ExpandPath expands a leading ~/ to the user home directory.
func ExpandPath(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("empty path")
	}

	if path == "~" || strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		if path == "~" {
			return home, nil
		}

		return filepath.Join(home, path[2:]), nil
	}

	return filepath.Clean(path), nil
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
	if err := os.MkdirAll(filepath.Dir(filename), 0o755); err != nil { //nolint:gosec
		return err
	}

	if c.legacy {
		return c.saveLegacy(filename)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0o644) //nolint:gosec
}

// Load profiles from json file
func (c *Config) Load(filename string) (err error) {
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		// New files use the map format, except explicit legacy paths.
		c.legacy = isLegacyPath(filename)

		return c.Save(filename)
	}

	body, err := os.ReadFile(filename) //nolint:gosec
	if err != nil {
		return err
	}

	legacy, err := detectLegacyFormat(body, filename)
	if err != nil {
		return err
	}

	c.legacy = legacy
	c.Profiles = make(map[string]Entry)

	if legacy {
		oldCfg := oldconfig.New()

		err = json.Unmarshal(body, oldCfg)
		if err != nil {
			return err
		}

		for profileName, oldEntries := range oldCfg.Profiles {
			for _, entry := range oldEntries {
				c.Store(profileName, entry.Key, entry.Value)
			}
		}

		return nil
	}

	err = json.Unmarshal(body, c)
	if err != nil {
		return err
	}

	if c.Profiles == nil {
		c.Profiles = make(map[string]Entry)
	}

	return nil
}

func (c *Config) saveLegacy(filename string) error {
	oldCfg := oldconfig.New()

	for profileName, entries := range c.Profiles {
		keys := make([]string, 0, len(entries))
		for key := range entries {
			keys = append(keys, key)
		}

		sort.Strings(keys)

		for _, key := range keys {
			oldCfg.Profiles[profileName] = append(oldCfg.Profiles[profileName], oldconfig.OldEntry{
				Key:   key,
				Value: entries[key],
			})
		}
	}

	return oldCfg.Save(filename)
}

func xdgConfigPath() (string, error) {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		configHome = filepath.Join(home, ".config")
	}

	return filepath.Join(configHome, xdgAppDir, xdgConfigFile), nil
}

func legacyConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, legacyFileName), nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}

func isLegacyPath(path string) bool {
	return filepath.Base(path) == legacyFileName
}

// detectLegacyFormat inspects profile values to choose array vs map encoding.
// Empty configs keep the legacy format when loaded from ~/.gitprofile.
func detectLegacyFormat(body []byte, filename string) (bool, error) {
	var probe struct {
		Profiles map[string]json.RawMessage `json:"profiles"`
	}

	err := json.Unmarshal(body, &probe)
	if err != nil {
		return false, err
	}

	if len(probe.Profiles) == 0 {
		return isLegacyPath(filename), nil
	}

	for _, raw := range probe.Profiles {
		raw = trimJSONSpace(raw)
		if len(raw) == 0 {
			continue
		}

		switch raw[0] {
		case '[':
			return true, nil
		case '{':
			return false, nil
		}
	}

	return isLegacyPath(filename), nil
}

func trimJSONSpace(raw json.RawMessage) json.RawMessage {
	return json.RawMessage(strings.TrimSpace(string(raw)))
}
