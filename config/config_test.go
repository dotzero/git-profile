package config

import (
	"testing"

	"github.com/matryer/is"
)

func TestNewConfig(t *testing.T) {
	is := is.New(t)

	var expected Config
	expected.Profiles = make(map[string][]Entry)

	is.Equal(&expected, NewConfig())
}

func TestGetProfile(t *testing.T) {
	is := is.New(t)

	var c Config
	c.Profiles = make(map[string][]Entry)
	c.Profiles["foo"] = []Entry{}

	_, ok := c.GetProfile("foo")
	is.True(ok)

	_, ok = c.GetProfile("bar")
	is.True(!ok)

	delete(c.Profiles, "foo")
	_, ok = c.GetProfile("foo")
	is.True(!ok)
}

func TestRemoveProfile(t *testing.T) {
	is := is.New(t)

	var c Config
	c.Profiles = make(map[string][]Entry)
	c.Profiles["foo"] = []Entry{}

	ok := c.RemoveProfile("foo")
	is.True(ok)

	ok = c.RemoveProfile("bar")
	is.True(ok)
}

func TestSetValue(t *testing.T) {
	is := is.New(t)

	data := []struct {
		Profile  string
		Key      string
		Value    string
		Expected *Config
	}{
		{
			Profile: "foo",
			Key:     "key1",
			Value:   "value1",
			Expected: &Config{
				Profiles: map[string][]Entry{
					"foo": []Entry{
						{
							Key:   "key1",
							Value: "value1",
						},
					},
				},
			},
		},
		{
			Profile: "foo",
			Key:     "key1",
			Value:   "value2",
			Expected: &Config{
				Profiles: map[string][]Entry{
					"foo": []Entry{
						{
							Key:   "key1",
							Value: "value2",
						},
					},
				},
			},
		},
		{
			Profile: "foo",
			Key:     "key2",
			Value:   "value2",
			Expected: &Config{
				Profiles: map[string][]Entry{
					"foo": []Entry{
						{
							Key:   "key1",
							Value: "value2",
						},
						{
							Key:   "key2",
							Value: "value2",
						},
					},
				},
			},
		},
		{
			Profile: "bar",
			Key:     "key1",
			Value:   "value1",
			Expected: &Config{
				Profiles: map[string][]Entry{
					"foo": []Entry{
						{
							Key:   "key1",
							Value: "value2",
						},
						{
							Key:   "key2",
							Value: "value2",
						},
					},
					"bar": []Entry{
						{
							Key:   "key1",
							Value: "value1",
						},
					},
				},
			},
		},
	}

	c := NewConfig()
	for _, tc := range data {
		c.SetValue(tc.Profile, Entry{Key: tc.Key, Value: tc.Value})
		is.Equal(tc.Expected, c)
	}
}

func TestRemoveValue(t *testing.T) {
	is := is.New(t)

	c := NewConfig()
	c.SetValue("foo", Entry{Key: "key1", Value: "value1"})
	c.SetValue("foo", Entry{Key: "key2", Value: "value2"})
	c.RemoveValue("foo", "key1")

	expected := &Config{
		Profiles: map[string][]Entry{
			"foo": []Entry{
				{
					Key:   "key2",
					Value: "value2",
				},
			},
		},
	}

	is.Equal(expected, c)
}
