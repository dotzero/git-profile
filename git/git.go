package git

import (
	"os/exec"
	"strings"
)

// Git is a vcs
type Git struct{}

// IsRepository checks that current directory is a git repository
func (g *Git) IsRepository() bool {
	err := exec.Command("git", "rev-parse", "--git-dir").Run()

	return err == nil
}

// Get returns the value stored in git local config
func (g *Git) Get(key string) (string, error) {
	out, err := exec.Command("git", "config", "--local", key).Output()

	return strings.TrimSpace(string(out)), err
}

// Set sets the value for a key in git local config
func (g *Git) Set(key string, value string) error {
	return exec.Command("git", "config", "--local", key, value).Run()
}
