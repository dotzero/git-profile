package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/matryer/is"

	"github.com/dotzero/git-profile/internal/config"
)

func TestMigrateCommand(t *testing.T) {
	is := is.New(t)

	home := t.TempDir()
	xdgHome := filepath.Join(home, "xdg")
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", xdgHome)

	legacyPath := filepath.Join(home, ".gitprofile")
	err := os.WriteFile(legacyPath, []byte(`{
  "profiles": {
    "work": [
      {"key": "user.email", "value": "work@example.com"}
    ]
  }
}`), 0o644)
	is.NoErr(err)

	var b bytes.Buffer

	cmd := Migrate()
	cmd.SetOut(&b)
	cmd.SetErr(&b)
	err = cmd.Execute()
	is.NoErr(err)

	xdgPath := filepath.Join(xdgHome, "git-profile", "config.json")
	body, err := os.ReadFile(xdgPath)
	is.NoErr(err)

	var decoded config.Config

	err = json.Unmarshal(body, &decoded)
	is.NoErr(err)
	is.Equal(decoded.Profiles["work"]["user.email"], "work@example.com")
}
