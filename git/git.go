package git

import (
	"os"
	"os/exec"
)

// IsRepository check that current directory is a root of git repository
func IsRepository() bool {
	s, err := os.Stat(".git")
	return err == nil && s.IsDir()
}

// SetConfig set git config key with value
func SetConfig(key string, value string) error {
	if err := exec.Command("git", "config", "--local", key, value).Run(); err != nil {
		return err
	}
	return nil
}
