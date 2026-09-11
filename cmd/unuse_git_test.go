package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestUnuseWithGitValues(t *testing.T) {
	for _, values := range [][]string{nil, {"work@example.com"}, {"work@example.com", "work@example.com", "extra@example.com"}} {
		t.Run(fmt.Sprintf("%d values", len(values)), func(t *testing.T) {
			is := is.New(t)
			repo, filename := setupCLIExitTest(t, false)
			_, err := exec.Command("git", "-C", repo, "config", "--local", "--unset", "user.email").CombinedOutput()
			is.NoErr(err)

			for _, value := range values {
				_, err := exec.Command("git", "-C", repo, "config", "--local", "--add", "user.email", value).CombinedOutput()
				is.NoErr(err)
			}

			globalPath := filepath.Join(t.TempDir(), "gitconfig")
			globalConfig := []byte("[user]\n\temail = global@example.com\n\temail = other@example.com\n")
			is.NoErr(os.WriteFile(globalPath, globalConfig, 0o600))
			t.Setenv("GIT_CONFIG_GLOBAL", globalPath)

			// The second invocation must also succeed after the local keys are gone.
			for _, args := range [][]string{{"unuse"}, {"unuse", "work"}} {
				child := cliExitTestCommand(repo, filename, args...)

				var stdout, stderr bytes.Buffer

				child.Stdout = &stdout
				child.Stderr = &stderr
				is.NoErr(child.Run())
				is.Equal(child.ProcessState.ExitCode(), 0)
				is.True(strings.Contains(stdout.String(), "Successfully removed"))
				is.Equal(stderr.String(), "")

				for _, key := range []string{"user.email", currentProfileKey} {
					query := exec.Command("git", "-C", repo, "config", "--local", "--get-all", key)
					out, err := query.Output()
					is.True(err != nil)
					is.Equal(query.ProcessState.ExitCode(), 1)
					is.Equal(string(out), "")
				}
			}

			afterGlobal, err := os.ReadFile(globalPath)
			is.NoErr(err)
			is.Equal(afterGlobal, globalConfig)
		})
	}
}

func TestUnuseReportsGitErrors(t *testing.T) {
	for _, failure := range []string{"locked config", "invalid key", "unwritable directory"} {
		t.Run(failure, func(t *testing.T) {
			is := is.New(t)
			repo, filename := setupCLIExitTest(t, false)
			gitDir := filepath.Join(repo, ".git")
			gitConfig := filepath.Join(gitDir, "config")
			before, err := os.ReadFile(gitConfig)
			is.NoErr(err)

			switch failure {
			case "locked config":
				is.NoErr(os.WriteFile(gitConfig+".lock", []byte("locked"), 0o600))
			case "invalid key":
				is.NoErr(os.WriteFile(filename, []byte(`{"profiles":{"work":{"invalid":"value"}}}`), 0o600))
			case "unwritable directory":
				if os.Geteuid() == 0 {
					t.Skip("root can write to directories regardless of mode bits")
				}

				info, err := os.Stat(gitDir)
				is.NoErr(err)
				t.Cleanup(func() { is.NoErr(os.Chmod(gitDir, info.Mode().Perm())) })
				is.NoErr(os.Chmod(gitDir, 0o500))
			}

			child := cliExitTestCommand(repo, filename, "unuse", "work")

			var stdout, stderr bytes.Buffer

			child.Stdout = &stdout
			child.Stderr = &stderr
			is.True(child.Run() != nil)
			is.Equal(child.ProcessState.ExitCode(), 1)
			is.Equal(stdout.String(), "")
			is.True(strings.Contains(stderr.String(), "Unable to interact with git to remove current profile"))

			after, err := os.ReadFile(gitConfig)
			is.NoErr(err)
			is.Equal(after, before)
		})
	}
}
