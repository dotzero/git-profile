package git

import (
	"testing"

	"github.com/matryer/is"
)

func TestIsRepository(t *testing.T) {
	is := is.New(t)

	ok := IsRepository()
	is.True(ok)
}

func TestLead(t *testing.T) {
	is := is.New(t)

	out, err := Lead("core.bare")
	is.NoErr(err)
	is.Equal(out, "false")
}

func TestStore(t *testing.T) {
	is := is.New(t)

	err := Store("test.test", "true")
	is.NoErr(err)
}
