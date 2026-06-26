package cmd

import (
	"bytes"
	"testing"

	"github.com/matryer/is"

	"github.com/dotzero/git-profile/internal/config"
)

func TestList(t *testing.T) {
	is := is.New(t)

	cfg := &storageMock{
		LenFunc: func() int {
			return 1
		},
		NamesFunc: func() []string {
			return []string{"home"}
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
		GetFunc: func(key string) (string, error) {
			return "home", nil
		},
	}

	var b bytes.Buffer

	cmd := List(cfg, vcs)

	cmd.SetOut(&b)
	err := cmd.Execute()

	is.NoErr(err)
	is.Equal(trim(b.String()), "Available profiles:\n- home:\n  user.email: work@example.com")
}
