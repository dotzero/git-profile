package cmd

import (
	"bytes"
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
	cmd.SetArgs([]string{"profile", `{"user.email": "work@example.com"}`})
	err := cmd.Execute()

	is.NoErr(err)
	is.Equal(trim(b.String()), "Successfully imported `profile` profile.")
}
