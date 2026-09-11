package cmd

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/matryer/is"

	"github.com/dotzero/git-profile/internal/config"
)

func TestList(t *testing.T) {
	is := is.New(t)

	cfg := config.New()
	cfg.Store("work", "user.email", "work@example.com")
	cfg.Store("home", "user.email", "home@example.com")
	cfg.Store("home", "core.autocrlf", "input")

	var b bytes.Buffer

	cmd := List(cfg)

	cmd.SetOut(&b)
	err := cmd.Execute()

	is.NoErr(err)
	is.Equal(trim(b.String()), multiline(
		"Available profiles:",
		"- home:",
		"  core.autocrlf: input",
		"  user.email: home@example.com",
		"- work:",
		"  user.email: work@example.com",
	))
}

func TestListOutsideRepository(t *testing.T) {
	for _, empty := range []bool{false, true} {
		t.Run(fmt.Sprintf("empty=%t", empty), func(t *testing.T) {
			is := is.New(t)
			repo, filename := setupCLIExitTest(t, empty)
			inside := cliExitTestCommand(repo, filename, "list")

			var insideOut, insideErr bytes.Buffer

			inside.Stdout = &insideOut
			inside.Stderr = &insideErr
			is.NoErr(inside.Run())

			// Listing profiles must also work when Git is unavailable.
			t.Setenv("PATH", t.TempDir())
			outsideDir := t.TempDir()
			outside := cliExitTestCommand(outsideDir, filename, "list")
			outside.Dir = outsideDir

			var outsideOut, outsideErr bytes.Buffer

			outside.Stdout = &outsideOut
			outside.Stderr = &outsideErr
			is.NoErr(outside.Run())
			is.Equal(outside.ProcessState.ExitCode(), 0)
			is.Equal(outsideOut.String(), insideOut.String())
			is.Equal(outsideErr.String(), insideErr.String())

			if empty {
				is.True(strings.Contains(outsideErr.String(), "There are no available profiles."))
			} else {
				is.Equal(outsideErr.String(), "")
				is.True(strings.Contains(outsideOut.String(), "work@example.com"))
			}
		})
	}
}
