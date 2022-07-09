package git

import (
	"os/exec"
	"strings"
)

// Git is a vcs
type Git struct {
	exec func(name string, arg ...string) *exec.Cmd
}

// New initializes and returns a new Git
func New() *Git {
	return &Git{
		exec: exec.Command,
	}
}

// IsRepository checks that current directory is a git repository
func (g *Git) IsRepository() bool {
	_, err := g.exec("git", "rev-parse", "--git-dir").CombinedOutput()

	return err == nil
}

// Get returns the value stored in git local config
func (g *Git) Get(key string) (string, error) {
	out, err := g.exec("git", "config", "--local", key).CombinedOutput()

	return strings.TrimSpace(string(out)), err
}

// Set sets the value for a key in git local config
func (g *Git) Set(key string, value string) error {
	_, err := g.exec("git", "config", "--local", key, value).CombinedOutput()

	return err
}
