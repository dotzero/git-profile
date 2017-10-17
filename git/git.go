package git

import (
	"log"
	"os/exec"
)

// IsRepository check that current directory is a git repository
func IsRepository() bool {
	log.Println("[DEBUG] IsRepository")
	err := exec.Command("git", "rev-parse", "--git-dir").Run()
	if err != nil {
		return false
	}
	return true
}

// SetLocalConfig set git local config key with value
func SetLocalConfig(key string, value string) error {
	log.Printf("[DEBUG] git config --local %s \"%s\"\n", key, value)
	if err := exec.Command("git", "config", "--local", key, value).Run(); err != nil {
		return err
	}
	return nil
}

// GetLocalConfig get git local config value by key
func GetLocalConfig(key string) ([]byte, error) {
	log.Printf("[DEBUG] git config --local %s \n", key)
	return exec.Command("git", "config", "--local", key).Output()
}
