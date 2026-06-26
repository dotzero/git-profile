package cmd

import (
	"bytes"
	"testing"

	"github.com/matryer/is"

	"github.com/dotzero/git-profile/internal/config"
)

func TestExport(t *testing.T) {
	is := is.New(t)

	cfg := &storageMock{
		LookupFunc: func(name string) (config.Entry, bool) {
			return config.Entry{
				"user.email": "work@example.com",
			}, true
		},
	}

	var b bytes.Buffer

	cmd := Export(cfg)

	cmd.SetOut(&b)
	cmd.SetArgs([]string{"profile"})
	err := cmd.Execute()

	is.NoErr(err)
	is.Equal(trim(b.String()), `{"user.email":"work@example.com"}`)
}
