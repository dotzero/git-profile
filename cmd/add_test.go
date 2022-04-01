package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestAdd(t *testing.T) {
	is := is.New(t)

	cfg := &storageMock{
		SaveFunc: func(filename string) error {
			return nil
		},
		StoreFunc: func(profile string, key string, value string) {},
	}

	var b bytes.Buffer

	cmd := Add(cfg)

	cmd.SetOut(&b)
	cmd.SetArgs([]string{"profile", "key", "value"})
	err := cmd.Execute()

	is.NoErr(err)
	is.Equal(strings.TrimSpace(b.String()), "Successfully stored `key=value` to `profile` profile.")
}
