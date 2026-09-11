package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestRegisterFlagsAddsDirectoryFlag(t *testing.T) {
	is := is.New(t)

	cmd := New()
	cmd.registerFlags()

	flag := cmd.PersistentFlags().Lookup("directory")

	is.True(flag != nil)
	is.Equal(flag.Shorthand, "C")
}

func TestServiceCommandsSkipConfig(t *testing.T) {
	commands := [][]string{
		{"version"},
		{"completion", "bash"},
		{"completion", "zsh"},
		{"completion", "fish"},
		{"completion", "powershell"},
	}

	for _, args := range commands {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			is := is.New(t)
			dir := t.TempDir()
			validConfig := filepath.Join(dir, "valid.json")
			is.NoErr(os.WriteFile(validConfig, []byte(`{"profiles":{}}`), 0o600))

			baseline := cliExitTestCommand(dir, validConfig, args...)
			wantOutput, err := baseline.Output()
			is.NoErr(err)
			is.True(len(wantOutput) > 0)

			for _, state := range []string{"broken", "missing"} {
				t.Run(state, func(t *testing.T) {
					is := is.New(t)
					filename := filepath.Join(t.TempDir(), "profiles.json")
					brokenConfig := []byte("{broken")

					if state == "broken" {
						is.NoErr(os.WriteFile(filename, brokenConfig, 0o600))
					}

					child := cliExitTestCommand(dir, filename, args...)
					child.Dir = dir

					var stderr bytes.Buffer

					child.Stderr = &stderr
					output, err := child.Output()
					is.NoErr(err)
					is.Equal(child.ProcessState.ExitCode(), 0)
					is.Equal(output, wantOutput)
					is.Equal(stderr.String(), "")

					after, err := os.ReadFile(filename)
					if state == "broken" {
						is.NoErr(err)
						is.Equal(after, brokenConfig)
					} else {
						is.True(os.IsNotExist(err))
					}
				})
			}
		})
	}
}

func TestProfileCommandsRejectBrokenConfig(t *testing.T) {
	commands := [][]string{
		{"list"},
		{"current"},
		{"use", "work"},
		{"unuse", "work"},
		{"export", "work"},
		{"import", "work", `{}`},
		{"add", "work", "user.email", "work@example.com"},
		{"del", "work"},
	}

	for _, args := range commands {
		t.Run(args[0], func(t *testing.T) {
			is := is.New(t)
			dir := t.TempDir()
			filename := filepath.Join(dir, "profiles.json")
			brokenConfig := []byte("{broken")
			is.NoErr(os.WriteFile(filename, brokenConfig, 0o600))

			child := cliExitTestCommand(dir, filename, args...)

			var stdout, stderr bytes.Buffer

			child.Stdout = &stdout
			child.Stderr = &stderr
			is.True(child.Run() != nil)
			is.Equal(child.ProcessState.ExitCode(), 1)
			is.True(strings.Contains(stderr.String(), "Unable to load config:"))
			is.Equal(stdout.String(), "")

			after, err := os.ReadFile(filename)
			is.NoErr(err)
			is.Equal(after, brokenConfig)
		})
	}
}
