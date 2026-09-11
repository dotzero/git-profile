package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/matryer/is"
)

// TestCLIProcess runs the application's entry point in a subprocess so os.Exit
// can be tested without terminating the parent test process.
func TestCLIProcess(t *testing.T) {
	if os.Getenv("GIT_PROFILE_TEST_CLI") != "1" {
		return
	}

	separator := slices.Index(os.Args, "--")
	if separator < 0 {
		t.Fatal("missing CLI argument separator")
	}

	app := New()
	app.SetArgs(os.Args[separator+1:])
	app.Execute()
	os.Exit(0)
}

func TestMissingProfileExitCodes(t *testing.T) {
	for _, empty := range []bool{false, true} {
		t.Run(fmt.Sprintf("empty=%t", empty), func(t *testing.T) {
			is := is.New(t)
			repo, filename := setupCLIExitTest(t, empty)
			gitConfig := filepath.Join(repo, ".git", "config")
			beforeGit, err := os.ReadFile(gitConfig)
			is.NoErr(err)
			beforeProfiles, err := os.ReadFile(filename)
			is.NoErr(err)

			for _, command := range []string{"use", "export", "unuse"} {
				t.Run(command, func(t *testing.T) {
					is := is.New(t)
					child := cliExitTestCommand(repo, filename, command, "typo")

					var stdout, stderr bytes.Buffer

					child.Stdout = &stdout
					child.Stderr = &stderr
					err := child.Run()
					is.True(err != nil)
					is.Equal(child.ProcessState.ExitCode(), 1)
					is.Equal(stdout.String(), "")
					is.True(strings.Contains(stderr.String(), "There is no profile with `typo` name"))

					afterGit, err := os.ReadFile(gitConfig)
					is.NoErr(err)
					is.Equal(afterGit, beforeGit)

					afterProfiles, err := os.ReadFile(filename)
					is.NoErr(err)
					is.Equal(afterProfiles, beforeProfiles)
				})
			}

			t.Run("shell chain stops", func(t *testing.T) {
				is := is.New(t)
				cli := cliExitTestCommand(repo, filename, "use", "typo")
				args := append([]string{"-c", `"$@" && printf continued`, "git-profile-test"}, cli.Args...)
				child := exec.Command("sh", args...)
				child.Env = cli.Env

				var stdout, stderr bytes.Buffer

				child.Stdout = &stdout
				child.Stderr = &stderr
				err := child.Run()
				is.True(err != nil)
				is.Equal(child.ProcessState.ExitCode(), 1)
				is.Equal(stdout.String(), "")
			})
		})
	}
}

func TestSuccessfulCLIExitCodes(t *testing.T) {
	for _, command := range []string{"use", "export", "unuse"} {
		t.Run(command, func(t *testing.T) {
			is := is.New(t)
			repo, filename := setupCLIExitTest(t, false)
			child := cliExitTestCommand(repo, filename, command, "work")
			out, err := child.CombinedOutput()
			is.NoErr(err)
			is.Equal(child.ProcessState.ExitCode(), 0)
			is.True(len(out) > 0)
		})
	}
}

func TestUnuseWithoutCurrentProfileExitCode(t *testing.T) {
	for _, empty := range []bool{false, true} {
		t.Run(fmt.Sprintf("empty=%t", empty), func(t *testing.T) {
			is := is.New(t)
			repo, filename := setupCLIExitTest(t, empty)
			_, err := exec.Command("git", "-C", repo, "config", "--local", "--unset", currentProfileKey).CombinedOutput()
			is.NoErr(err)

			child := cliExitTestCommand(repo, filename, "unuse")
			_, err = child.CombinedOutput()
			is.NoErr(err)
			is.Equal(child.ProcessState.ExitCode(), 0)
		})
	}
}

func cliExitTestCommand(repo, filename string, args ...string) *exec.Cmd {
	childArgs := append([]string{"-test.run=^TestCLIProcess$", "--", "--config", filename, "-C", repo}, args...)
	child := exec.Command(os.Args[0], childArgs...)

	child.Env = append(os.Environ(), "GIT_PROFILE_TEST_CLI=1")

	return child
}

func setupCLIExitTest(t *testing.T, empty bool) (string, string) {
	t.Helper()
	is := is.New(t)
	repo := t.TempDir()
	filename := filepath.Join(t.TempDir(), "profiles.json")
	t.Setenv("GIT_CONFIG_GLOBAL", os.DevNull)
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	t.Setenv("GIT_CONFIG_COUNT", "0")

	for _, key := range []string{"GIT_DIR", "GIT_WORK_TREE", "GIT_COMMON_DIR", "GIT_CONFIG", "GIT_CONFIG_PARAMETERS"} {
		t.Setenv(key, "")
		is.NoErr(os.Unsetenv(key))
	}

	out, err := exec.Command("git", "-c", "init.templateDir=", "init", "-q", repo).CombinedOutput()
	if err != nil {
		t.Fatalf("git init: %s: %v", out, err)
	}

	for key, value := range map[string]string{currentProfileKey: "work", "user.email": "work@example.com"} {
		_, err := exec.Command("git", "-C", repo, "config", "--local", key, value).CombinedOutput()
		is.NoErr(err)
	}

	profiles := `{"profiles":{"work":{"user.email":"work@example.com"}}}`
	if empty {
		profiles = `{"profiles":{}}`
	}

	is.NoErr(os.WriteFile(filename, []byte(profiles), 0o600))

	return repo, filename
}
