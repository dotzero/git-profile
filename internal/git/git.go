package git

import (
	"errors"
	"os/exec"
	"strings"
)

// gitConfigKeyNotFound is the exit code git returns from
// `git config --unset` when the key is not present.
const gitConfigKeyNotFound = 5

// Git is a vcs
type Git struct {
	exec func(name string, arg ...string) *exec.Cmd
	dir  string
}

// New initializes and returns a new Git
func New() *Git {
	return &Git{
		exec: exec.Command,
	}
}

// SetDir sets the directory where git commands will be executed.
func (g *Git) SetDir(dir string) {
	g.dir = dir
}

// IsRepository checks that current directory is a git repository
func (g *Git) IsRepository() bool {
	_, err := g.command("rev-parse", "--git-dir").CombinedOutput()

	return err == nil
}

// Get returns the value stored in git local config
func (g *Git) Get(key string) (string, error) {
	out, err := g.command("config", "--local", key).CombinedOutput()

	return strings.TrimSpace(string(out)), err
}

// Set sets the value for a key in git local config
func (g *Git) Set(key string, value string) error {
	_, err := g.command("config", "--local", key, value).CombinedOutput()

	return err
}

// Unset removes a key from git local config.
// A missing key is not treated as an error.
func (g *Git) Unset(key string) error {
	_, err := g.command("config", "--local", "--unset", key).CombinedOutput()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() == gitConfigKeyNotFound {
			return nil
		}

		return err
	}

	return nil
}

func (g *Git) command(arg ...string) *exec.Cmd {
	cmd := g.exec("git", arg...)
	cmd.Dir = g.dir

	return cmd
}
