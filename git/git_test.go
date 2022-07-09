package git

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/matryer/is"
)

// This one trick you won't believe to mock an external binary within test!
//
// Creates an exec.Cmd that actually calls back into the test binary itself,
// invoking a specific test function, with [-- cmd, args...] appended to end,
// and a magic env var specified.
//
// This allows us to easily mock external binary behavior in a cross-platform
// way. The go stdlib uses this trick in os/exec.Cmd's own tests!
func mockExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)

	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}

	return cmd
}

func TestHelperProcess(t *testing.T) {
	// ignore me when not being specifically invoked
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	var args []string

	// grab all the args after "--"
	for i, arg := range os.Args {
		if arg == "--" {
			args = os.Args[i+1:]
			break
		}
	}

	if len(args) == 0 {
		t.Fatal("mock command not specified")
	}

	switch len(args) {
	case 3: // git rev-parse --git-dir
		fmt.Fprint(os.Stdout, ".git")
		os.Exit(0)
	case 4: // git config --local user.name
		fmt.Fprint(os.Stdout, "test_user")
		os.Exit(0)
	case 5: // git config --local user.name test
		os.Exit(0)
	}

	t.Fatal("mocked command not implemented:", args)
}

func TestIsRepository(t *testing.T) {
	is := is.New(t)

	g := &Git{exec: mockExecCommand}
	ok := g.IsRepository()

	is.True(ok)
}

func TestGet(t *testing.T) {
	is := is.New(t)

	g := &Git{exec: mockExecCommand}
	value, err := g.Get("user.name")

	is.NoErr(err)
	is.Equal(value, "test_user")
}

func TestSet(t *testing.T) {
	is := is.New(t)

	g := &Git{exec: mockExecCommand}
	err := g.Set("user.name", "test_user")

	is.NoErr(err)
}
