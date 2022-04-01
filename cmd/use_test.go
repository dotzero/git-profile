package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/matryer/is"

	"github.com/dotzero/git-profile/config"
)

func TestUse(t *testing.T) {
	is := is.New(t)

	cfg := &storageMock{
		LookupFunc: func(name string) ([]config.Entry, bool) {
			return []config.Entry{
				{Key: "user.email", Value: "work@example.com"},
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
	is.Equal(strings.TrimSpace(b.String()), "Successfully applied `profile` profile to current git repository.")
}
