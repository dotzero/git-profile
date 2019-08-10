package git

import (
	"os/exec"
	"strings"
)

// IsRepository checks that current directory is a git repository
func IsRepository() bool {
	err := exec.Command("git", "rev-parse", "--git-dir").Run()
	if err != nil { //nolint
		return false
	}
	return true
}

// Lead returns the value stored in git local config
func Lead(key string) (string, error) {
	out, err := exec.Command("git", "config", "--local", key).Output()
	return strings.TrimSpace(string(out)), err
}

// Store sets the value for a key in git local config
func Store(key string, value string) error {
	return exec.Command("git", "config", "--local", key, value).Run()
}
