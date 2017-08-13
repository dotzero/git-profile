package config

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

func TestNewConfig(t *testing.T) {
	var expected Config
	expected.Profiles = make(map[string][]Entry)

	equals(t, &expected, NewConfig())
}

func TestGetProfile(t *testing.T) {
	var c Config
	c.Profiles = make(map[string][]Entry)
	c.Profiles["foo"] = []Entry{}

	_, ok := c.GetProfile("foo")
	assert(t, ok, "expected true on exists profile")

	_, ok = c.GetProfile("bar")
	assert(t, !ok, "expected false on non exists profile")

	delete(c.Profiles, "foo")
	_, ok = c.GetProfile("foo")
	assert(t, !ok, "expected false on deleted profile")
}

func TestRemoveProfile(t *testing.T) {
	var c Config
	c.Profiles = make(map[string][]Entry)
	c.Profiles["foo"] = []Entry{}

	ok := c.RemoveProfile("foo")
	assert(t, ok, "expected true on exists profile")

	ok = c.RemoveProfile("bar")
	assert(t, ok, "expected true on non exists profile")
}

func TestSetValue(t *testing.T) {
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
		equals(t, tc.Expected, c)
	}
}

func TestRemoveValue(t *testing.T) {
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

	equals(t, expected, c)
}
