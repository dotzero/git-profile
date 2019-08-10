package config

import (
	"testing"

	"github.com/matryer/is"
)

func TestDelete(t *testing.T) {
	is := is.New(t)

	c := New()
	c.Store("foo", "key1", "value1")
	c.Store("foo", "key2", "value2")
	c.Delete("foo", "key1")

	expected := &Config{
		Profiles: map[string][]Entry{
			"foo": {
				{"key2", "value2"},
			},
		},
	}

	is.Equal(c, expected)
}

func TestSetValue(t *testing.T) {
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

	c := New()
	for _, tc := range cases {
		var tc = tc // pin
		c.Store(tc.profile, tc.key, tc.value)
		is.Equal(c, tc.expected)
	}
}
