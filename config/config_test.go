package config

import (
	"testing"

	"github.com/matryer/is"
)

func TestDelete(t *testing.T) {
	is := is.New(t)

	cfg := &Config{
		Profiles: map[string][]Entry{
			"home": {
				{"user.email", "work@example.com"},
				{"user.name", "John Doe"},
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
		Profiles: map[string][]Entry{
			"home": {
				{"user.email", "work@example.com"},
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
				Profiles: map[string][]Entry{
					"foo": {
						{"key1", "value1"},
					},
				},
			},
		},
		{
			profile: "foo",
			key:     "key1",
			value:   "value2",
			expected: &Config{
				Profiles: map[string][]Entry{
					"foo": {
						{"key1", "value2"},
					},
				},
			},
		},
		{
			profile: "foo",
			key:     "key2",
			value:   "value2",
			expected: &Config{
				Profiles: map[string][]Entry{
					"foo": {
						{"key1", "value2"},
						{"key2", "value2"},
					},
				},
			},
		},
		{
			profile: "bar",
			key:     "key1",
			value:   "value1",
			expected: &Config{
				Profiles: map[string][]Entry{
					"foo": {
						{"key1", "value2"},
						{"key2", "value2"},
					},
					"bar": {
						{"key1", "value1"},
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
