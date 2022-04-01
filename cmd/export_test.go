package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/matryer/is"

	"github.com/dotzero/git-profile/config"
)

func TestExport(t *testing.T) {
	is := is.New(t)

	cfg := &storageMock{
		LookupFunc: func(name string) ([]config.Entry, bool) {
			return []config.Entry{
				{Key: "user.email", Value: "work@example.com"},
			}, true
		},
	}

	var b bytes.Buffer

	cmd := Export(cfg)

	cmd.SetOut(&b)
	cmd.SetArgs([]string{"profile"})
	err := cmd.Execute()

	is.NoErr(err)
	is.Equal(strings.TrimSpace(b.String()), `[{"key":"user.email","value":"work@example.com"}]`)
}
