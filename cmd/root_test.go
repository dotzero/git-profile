package cmd

import (
	"testing"

	"github.com/matryer/is"
)

func TestRegisterFlagsAddsDirectoryFlag(t *testing.T) {
	is := is.New(t)

	cmd := New()
	cmd.registerFlags()

	flag := cmd.PersistentFlags().Lookup("directory")

	is.True(flag != nil)
	is.Equal(flag.Shorthand, "C")
}
