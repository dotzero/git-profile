package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestImport(t *testing.T) {
	is := is.New(t)

	cfg := &storageMock{
		SaveFunc: func(filename string) error {
			return nil
		},
		StoreFunc: func(profile string, key string, value string) {},
	}

	var b bytes.Buffer

	cmd := Import(cfg)

	cmd.SetOut(&b)
	cmd.SetArgs([]string{"profile", `[{"key":"user.email","value":"work@example.com"}]`})
	err := cmd.Execute()

	is.NoErr(err)
	is.Equal(strings.TrimSpace(b.String()), "Successfully import `profile` profile")
}
