package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestDelete(t *testing.T) {
	is := is.New(t)

	cfg := &Config{
		Profiles: map[string]Entry{
			"home": {
				"user.email": "work@example.com",
				"user.name":  "John Doe",
			},
		},
	}

	is.True(!cfg.Delete("work", "user.name"))
	is.Equal(len(cfg.Profiles), 1)
	is.Equal(len(cfg.Profiles["home"]), 2)
	is.True(cfg.Delete("home", "user.name"))
	is.Equal(len(cfg.Profiles), 1)
	is.Equal(len(cfg.Profiles["home"]), 1)
}

func TestDeleteProfile(t *testing.T) {
	is := is.New(t)

	cfg := &Config{
		Profiles: map[string]Entry{
			"home": {
				"user.email": "work@example.com",
			},
		},
	}

	is.True(!cfg.DeleteProfile("work"))
	is.Equal(len(cfg.Profiles), 1)
	is.True(cfg.DeleteProfile("home"))
	is.Equal(len(cfg.Profiles), 0)
}

func TestStoreValue(t *testing.T) {
	is := is.New(t)

	cases := []struct {
		profile  string
		key      string
		value    string
		expected *Config
	}{
		{
			profile: "foo",
			key:     "key1",
			value:   "value1",
			expected: &Config{
				Profiles: map[string]Entry{
					"foo": {
						"key1": "value1",
					},
				},
			},
		},
		{
			profile: "foo",
			key:     "key1",
			value:   "value2",
			expected: &Config{
				Profiles: map[string]Entry{
					"foo": {
						"key1": "value2",
					},
				},
			},
		},
		{
			profile: "foo",
			key:     "key2",
			value:   "value2",
			expected: &Config{
				Profiles: map[string]Entry{
					"foo": {
						"key1": "value2",
						"key2": "value2",
					},
				},
			},
		},
		{
			profile: "bar",
			key:     "key1",
			value:   "value1",
			expected: &Config{
				Profiles: map[string]Entry{
					"foo": {
						"key1": "value2",
						"key2": "value2",
					},
					"bar": {
						"key1": "value1",
					},
				},
			},
		},
	}

	cfg := New()

	for _, c := range cases {
		c := c // pin

		cfg.Store(c.profile, c.key, c.value)
		is.Equal(cfg.Profiles, c.expected.Profiles)
	}
}

func TestLoadLegacyDoesNotRewriteFile(t *testing.T) {
	is := is.New(t)

	dir := t.TempDir()
	path := filepath.Join(dir, ".gitprofile")
	original := `{
  "profiles": {
    "work": [
      {
        "key": "user.email",
        "value": "work@example.com"
      },
      {
        "key": "user.name",
        "value": "John Doe"
      }
    ]
  }
}`

	err := os.WriteFile(path, []byte(original), 0o644)
	is.NoErr(err)

	cfg := New()
	err = cfg.Load(path)
	is.NoErr(err)
	is.Equal(cfg.Profiles["work"]["user.email"], "work@example.com")
	is.Equal(cfg.Profiles["work"]["user.name"], "John Doe")
	is.True(cfg.legacy)

	body, err := os.ReadFile(path)
	is.NoErr(err)
	is.Equal(string(body), original)
}

func TestSaveLegacyKeepsArrayFormat(t *testing.T) {
	is := is.New(t)

	dir := t.TempDir()
	path := filepath.Join(dir, ".gitprofile")
	original := `{
  "profiles": {
    "work": [
      {
        "key": "user.email",
        "value": "work@example.com"
      }
    ]
  }
}`

	err := os.WriteFile(path, []byte(original), 0o644)
	is.NoErr(err)

	cfg := New()
	err = cfg.Load(path)
	is.NoErr(err)

	cfg.Store("work", "user.name", "John Doe")
	err = cfg.Save(path)
	is.NoErr(err)

	body, err := os.ReadFile(path)
	is.NoErr(err)

	var decoded map[string]map[string][]map[string]string

	err = json.Unmarshal(body, &decoded)
	is.NoErr(err)
	is.Equal(len(decoded["profiles"]["work"]), 2)

	keys := map[string]string{}
	for _, entry := range decoded["profiles"]["work"] {
		keys[entry["key"]] = entry["value"]
	}

	is.Equal(keys["user.email"], "work@example.com")
	is.Equal(keys["user.name"], "John Doe")
	is.True(!strings.Contains(string(body), `"user.email": "work@example.com"`))
}

func TestLoadEmptyLegacyPathKeepsArrayFormat(t *testing.T) {
	is := is.New(t)

	dir := t.TempDir()
	path := filepath.Join(dir, ".gitprofile")
	original := "{\n  \"profiles\": {}\n}"

	err := os.WriteFile(path, []byte(original), 0o644)
	is.NoErr(err)

	cfg := New()
	err = cfg.Load(path)
	is.NoErr(err)
	is.True(cfg.legacy)

	cfg.Store("work", "user.email", "work@example.com")
	err = cfg.Save(path)
	is.NoErr(err)

	body, err := os.ReadFile(path)
	is.NoErr(err)
	is.True(strings.Contains(string(body), `"key": "user.email"`))
	is.True(strings.Contains(string(body), `"value": "work@example.com"`))
}

func TestLoadSaveMapFormat(t *testing.T) {
	is := is.New(t)

	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	original := `{
  "profiles": {
    "work": {
      "user.email": "work@example.com"
    }
  }
}`

	err := os.WriteFile(path, []byte(original), 0o644)
	is.NoErr(err)

	cfg := New()
	err = cfg.Load(path)
	is.NoErr(err)
	is.True(!cfg.legacy)

	cfg.Store("work", "core.autocrlf", "input")
	err = cfg.Save(path)
	is.NoErr(err)

	body, err := os.ReadFile(path)
	is.NoErr(err)

	var decoded Config

	err = json.Unmarshal(body, &decoded)
	is.NoErr(err)
	is.Equal(decoded.Profiles["work"]["user.email"], "work@example.com")
	is.Equal(decoded.Profiles["work"]["core.autocrlf"], "input")
}

func TestDefaultPathPrefersXDG(t *testing.T) {
	is := is.New(t)

	home := t.TempDir()
	xdgHome := filepath.Join(home, "xdg")
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", xdgHome)

	xdgPath := filepath.Join(xdgHome, "git-profile", "config.json")
	legacyPath := filepath.Join(home, ".gitprofile")

	err := os.MkdirAll(filepath.Dir(xdgPath), 0o755)
	is.NoErr(err)
	err = os.WriteFile(xdgPath, []byte(`{"profiles":{}}`), 0o644)
	is.NoErr(err)
	err = os.WriteFile(legacyPath, []byte(`{"profiles":{}}`), 0o644)
	is.NoErr(err)

	path, err := DefaultPath()
	is.NoErr(err)
	is.Equal(path, xdgPath)
}

func TestDefaultPathFallsBackToLegacy(t *testing.T) {
	is := is.New(t)

	home := t.TempDir()
	xdgHome := filepath.Join(home, "xdg")
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", xdgHome)

	legacyPath := filepath.Join(home, ".gitprofile")
	err := os.WriteFile(legacyPath, []byte(`{"profiles":{}}`), 0o644)
	is.NoErr(err)

	path, err := DefaultPath()
	is.NoErr(err)
	is.Equal(path, legacyPath)
}

func TestDefaultPathUsesXDGWhenMissing(t *testing.T) {
	is := is.New(t)

	home := t.TempDir()
	xdgHome := filepath.Join(home, "xdg")
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", xdgHome)

	path, err := DefaultPath()
	is.NoErr(err)
	is.Equal(path, filepath.Join(xdgHome, "git-profile", "config.json"))
}

func TestExpandPath(t *testing.T) {
	is := is.New(t)

	home := t.TempDir()
	t.Setenv("HOME", home)

	path, err := ExpandPath("~/profiles.json")
	is.NoErr(err)
	is.Equal(path, filepath.Join(home, "profiles.json"))

	path, err = ExpandPath("~")
	is.NoErr(err)
	is.Equal(path, home)
}
