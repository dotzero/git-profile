package cmd

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"charm.land/huh/v2"
	"github.com/matryer/is"

	"github.com/dotzero/git-profile/internal/config"
	"github.com/dotzero/git-profile/internal/ui"
)

func TestAdd(t *testing.T) {
	is := is.New(t)

	cfg := &storageMock{
		SaveFunc: func(filename string) error {
			return nil
		},
		StoreFunc: func(profile string, key string, value string) {},
	}

	var b bytes.Buffer

	cmd := Add(cfg)

	cmd.SetOut(&b)
	cmd.SetArgs([]string{"profile", "key", "value"})
	err := cmd.Execute()

	is.NoErr(err)
	is.Equal(trim(b.String()), "Successfully stored `key=value` to `profile` profile.")
}

func TestAddRejectsPartialArgs(t *testing.T) {
	is := is.New(t)

	cmd := Add(&storageMock{})
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"profile"})

	err := cmd.Execute()

	is.True(err != nil)
	is.True(strings.Contains(err.Error(), "accepts either no arguments for interactive mode or exactly 3 arguments"))
}

func TestAddInteractiveStoresOnlyChangedFilledValues(t *testing.T) {
	is := is.New(t)

	cfg := &storageMock{
		LookupFunc: func(name string) (config.Entry, bool) {
			return config.Entry{
				userNameKey:  "Jane Doe",
				userEmailKey: "old@example.com",
			}, true
		},
		SaveFunc: func(filename string) error {
			return nil
		},
		StoreFunc: func(profile string, key string, value string) {},
	}

	var b bytes.Buffer

	cmd := addCommand(
		cfg,
		func(_ io.Reader, _ io.Writer) (string, error) {
			return "work", nil
		},
		func(initial ui.ProfileFormData, _ io.Reader, _ io.Writer) (ui.ProfileFormData, error) {
			is.Equal(initial.Profile, "work")
			is.Equal(initial.UserName, "Jane Doe")
			is.Equal(initial.UserEmail, "old@example.com")
			is.Equal(initial.UserSigningKey, "")

			return ui.ProfileFormData{
				Profile:        "work",
				UserName:       "Jane Doe",
				UserEmail:      "new@example.com",
				UserSigningKey: "",
			}, nil
		},
	)
	cmd.SetOut(&b)
	cmd.SetArgs(nil)

	err := cmd.Execute()

	is.NoErr(err)
	is.Equal(len(cfg.StoreCalls()), 1)
	is.Equal(cfg.StoreCalls()[0].Profile, "work")
	is.Equal(cfg.StoreCalls()[0].Key, userEmailKey)
	is.Equal(cfg.StoreCalls()[0].Value, "new@example.com")
	is.Equal(len(cfg.SaveCalls()), 1)
	is.Equal(trim(b.String()), "Successfully updated `work` profile.")
}

func TestAddInteractiveSkipsSaveWhenNothingChanged(t *testing.T) {
	is := is.New(t)

	cfg := &storageMock{
		LookupFunc: func(name string) (config.Entry, bool) {
			return config.Entry{
				userNameKey:  "Jane Doe",
				userEmailKey: "work@example.com",
			}, true
		},
		SaveFunc: func(filename string) error {
			return errors.New("save should not be called")
		},
		StoreFunc: func(profile string, key string, value string) {
			t.Fatalf("Store should not be called")
		},
	}

	var b bytes.Buffer

	cmd := addCommand(
		cfg,
		func(_ io.Reader, _ io.Writer) (string, error) {
			return "work", nil
		},
		func(initial ui.ProfileFormData, _ io.Reader, _ io.Writer) (ui.ProfileFormData, error) {
			is.Equal(initial.Profile, "work")

			return ui.ProfileFormData{
				Profile:        "work",
				UserName:       "Jane Doe",
				UserEmail:      "work@example.com",
				UserSigningKey: "",
			}, nil
		},
	)
	cmd.SetOut(&b)
	cmd.SetArgs(nil)

	err := cmd.Execute()

	is.NoErr(err)
	is.Equal(len(cfg.StoreCalls()), 0)
	is.Equal(len(cfg.SaveCalls()), 0)
	is.Equal(trim(b.String()), "No profile changes to save.")
}

func TestAddInteractiveAbortDoesNotPrintSaveError(t *testing.T) {
	is := is.New(t)

	cfg := &storageMock{
		StoreFunc: func(profile string, key string, value string) {
			t.Fatalf("Store should not be called")
		},
		SaveFunc: func(filename string) error {
			t.Fatalf("Save should not be called")
			return nil
		},
	}

	var (
		out    bytes.Buffer
		errOut bytes.Buffer
	)

	cmd := addCommand(
		cfg,
		func(_ io.Reader, _ io.Writer) (string, error) {
			return "", huh.ErrUserAborted
		},
		func(_ ui.ProfileFormData, _ io.Reader, _ io.Writer) (ui.ProfileFormData, error) {
			t.Fatalf("PromptProfileFields should not be called")
			return ui.ProfileFormData{}, nil
		},
	)
	cmd.SetOut(&out)
	cmd.SetErr(&errOut)
	cmd.SetArgs(nil)

	err := cmd.Execute()

	is.NoErr(err)
	is.Equal(trim(errOut.String()), "Interactive add cancelled.")
}
