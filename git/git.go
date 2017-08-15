package git

import (
	"log"
	"os"
	"os/exec"
)

// IsRepository check that current directory is a root of git repository
func IsRepository() bool {
	log.Println("[DEBUG] IsRepository")
	s, err := os.Stat(".git")
	return err == nil && s.IsDir()
}

// SetLocalConfig set git local config key with value
func SetLocalConfig(key string, value string) error {
	log.Printf("[DEBUG] git config --local %s \"%s\"\n", key, value)
	if err := exec.Command("git", "config", "--local", key, value).Run(); err != nil {
		return err
	}
	return nil
}
