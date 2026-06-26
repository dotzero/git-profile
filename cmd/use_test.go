package cmd

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/matryer/is"

	"github.com/dotzero/git-profile/internal/config"
)

func TestUse(t *testing.T) {
	is := is.New(t)

	cfg := &storageMock{
		LenFunc: func() int {
			return 1
		},
		LookupFunc: func(name string) (config.Entry, bool) {
			return config.Entry{
				"user.email": "work@example.com",
			}, true
		},
	}

	vcs := &vcsMock{
		IsRepositoryFunc: func() bool {
			return true
		},
		SetFunc: func(key string, value string) error {
			return nil
		},
	}

	var b bytes.Buffer

	cmd := Use(cfg, vcs)

	cmd.SetOut(&b)
	cmd.SetArgs([]string{"profile"})
	err := cmd.Execute()

	is.NoErr(err)
	is.Equal(trim(b.String()), "Successfully applied `profile` profile to current git repository.")
}

func TestProfileResolveInteractive(t *testing.T) {
	is := is.New(t)

	cfg := &storageMock{
		LenFunc: func() int {
			return 2
		},
		NamesFunc: func() []string {
			return []string{"home", "work"}
		},
	}

	profile, err := profileResolve(
		nil,
		&bytes.Buffer{},
		&bytes.Buffer{},
		cfg,
		func(names []string, _ io.Reader, _ io.Writer) (string, error) {
			is.Equal(names, []string{"home", "work"})
			return "work", nil
		},
	)

	is.NoErr(err)
	is.Equal(profile, "work")
}

func TestProfileResolveInteractiveError(t *testing.T) {
	is := is.New(t)

	cfg := &storageMock{
		LenFunc: func() int {
			return 1
		},
		NamesFunc: func() []string {
			return []string{"work"}
		},
	}

	_, err := profileResolve(
		nil,
		&bytes.Buffer{},
		&bytes.Buffer{},
		cfg,
		func(_ []string, _ io.Reader, _ io.Writer) (string, error) {
			return "", fmt.Errorf("boom")
		},
	)

	is.True(err != nil)
	is.Equal(err.Error(), "Unable to select a profile: boom")
}
