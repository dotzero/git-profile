package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/matryer/is"

	"github.com/dotzero/git-profile/config"
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
		LookupFunc: func(name string) ([]config.Entry, bool) {
			return []config.Entry{
				{Key: "user.email", Value: "work@example.com"},
			}, true
		},
	}

	var b bytes.Buffer

	cmd := List(cfg)

	cmd.SetOut(&b)
	err := cmd.Execute()

	is.NoErr(err)
	is.True(strings.Contains(b.String(), "home"))
	is.True(strings.Contains(b.String(), "user.email"))
}
