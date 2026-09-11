package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/matryer/is"

	"github.com/dotzero/git-profile/internal/config"
	"github.com/dotzero/git-profile/internal/git"
)

func TestUse(t *testing.T) {
	is := is.New(t)

	cfg := &storageMock{
		LenFunc: func() int {
			return 1
		},
		LookupFunc: func(name string) (config.Entry, bool) {
			return config.Entry{
				"user.email":    "work@example.com",
				"core.autocrlf": "input",
			}, true
		},
	}

	setParams := make(map[string]string)
	vcs := &vcsMock{
		IsRepositoryFunc: func() bool {
			return true
		},
		GetFunc: func(key string) (string, error) {
			return "", fmt.Errorf("no current profile")
		},
		SetFunc: func(key string, value string) error {
			setParams[key] = value
			return nil
		},
	}

	var b bytes.Buffer

	cmd := Use(cfg, vcs)

	cmd.SetOut(&b)
	cmd.SetArgs([]string{"profile"})
	err := cmd.Execute()

	is.NoErr(err)
	is.Equal(trim(b.String()), "Successfully applied `profile` profile to current git repository.")
	is.Equal(setParams, map[string]string{
		"current-profile.name": "profile",
		"user.email":           "work@example.com",
		"core.autocrlf":        "input",
	})
}

func TestProfileResolveInteractive(t *testing.T) {
	is := is.New(t)

	cfg := &storageMock{
		LenFunc: func() int {
			return 2
		},
		NamesFunc: func() []string {
			return []string{"home", "work"}
		},
	}

	profile, err := profileResolve(
		nil,
		&bytes.Buffer{},
		&bytes.Buffer{},
		cfg,
		func(names []string, _ io.Reader, _ io.Writer) (string, error) {
			is.Equal(names, []string{"home", "work"})
			return "work", nil
		},
	)

	is.NoErr(err)
	is.Equal(profile, "work")
}

func TestProfileResolveInteractiveError(t *testing.T) {
	is := is.New(t)

	cfg := &storageMock{
		LenFunc: func() int {
			return 1
		},
		NamesFunc: func() []string {
			return []string{"work"}
		},
	}

	_, err := profileResolve(
		nil,
		&bytes.Buffer{},
		&bytes.Buffer{},
		cfg,
		func(_ []string, _ io.Reader, _ io.Writer) (string, error) {
			return "", fmt.Errorf("boom")
		},
	)

	is.True(err != nil)
	is.Equal(err.Error(), "Unable to select a profile: boom")
}

func TestUseWithGit(t *testing.T) {
	tests := []struct {
		name     string
		previous string
		apply    []string
		wantMail string
		wantKey  string
	}{
		{
			name:  "switch removes previous profile keys",
			apply: []string{"work", "home"},
		},
		{
			name:  "reapplying profile succeeds",
			apply: []string{"work", "home", "home"},
		},
		{
			name:     "first application preserves existing settings",
			apply:    []string{"home"},
			wantMail: "existing@example.com",
			wantKey:  "EXISTING_KEY",
		},
		{
			name:     "unknown previous profile preserves existing settings",
			previous: "deleted-profile",
			apply:    []string{"home"},
			wantMail: "existing@example.com",
			wantKey:  "EXISTING_KEY",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			repo := t.TempDir()
			globalPath := filepath.Join(t.TempDir(), "gitconfig")
			globalConfig := []byte("[user]\n\temail = global@example.com\n\tsigningkey = GLOBAL_KEY\n")
			is.NoErr(os.WriteFile(globalPath, globalConfig, 0o600))

			t.Setenv("GIT_CONFIG_GLOBAL", globalPath)
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

			v := git.New()
			v.SetDir(repo)
			is.NoErr(v.Set("user.email", "existing@example.com"))
			is.NoErr(v.Set("user.signingkey", "EXISTING_KEY"))
			is.NoErr(v.Set("core.autocrlf", "input"))

			if tt.previous != "" {
				is.NoErr(v.Set(currentProfileKey, tt.previous))
			}

			cfg := config.New()
			cfg.Store("work", "user.name", "Work")
			cfg.Store("work", "user.email", "work@example.com")
			cfg.Store("work", "user.signingkey", "WORK_KEY")
			cfg.Store("home", "user.name", "Home")

			for _, profile := range tt.apply {
				cmd := Use(cfg, v)
				cmd.SetOut(io.Discard)
				cmd.SetArgs([]string{profile})
				is.NoErr(cmd.Execute())
			}

			for key, want := range map[string]string{
				currentProfileKey: "home",
				"user.name":       "Home",
				"user.email":      tt.wantMail,
				"user.signingkey": tt.wantKey,
				"core.autocrlf":   "input",
			} {
				got, err := v.Get(key)

				if want == "" {
					var exitErr *exec.ExitError
					if !errors.As(err, &exitErr) || exitErr.ExitCode() != 1 {
						t.Errorf("expected local %s to be absent, got %q, error %v", key, got, err)
					}
				} else {
					is.NoErr(err)
					is.Equal(got, want)
				}
			}

			gotGlobal, err := os.ReadFile(globalPath)
			is.NoErr(err)
			is.Equal(gotGlobal, globalConfig)
		})
	}
}
