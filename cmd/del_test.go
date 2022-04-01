package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestDel(t *testing.T) {
	is := is.New(t)

	cfg := &storageMock{
		DeleteProfileFunc: func(profile string) bool {
			return true
		},
		SaveFunc: func(filename string) error {
			return nil
		},
	}

	var b bytes.Buffer

	cmd := Del(cfg)

	cmd.SetOut(&b)
	cmd.SetArgs([]string{"profile"})
	err := cmd.Execute()

	is.NoErr(err)
	is.Equal(strings.TrimSpace(b.String()), "Successfully removed `profile` profile.")
}
