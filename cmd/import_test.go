package cmd

import (
	"bytes"
	"testing"

	"github.com/matryer/is"
)

func TestImport(t *testing.T) {
	t.Run("map format", func(t *testing.T) {
		is := is.New(t)

		stored := make(map[string]string)
		cfg := &storageMock{
			SaveFunc: func(filename string) error {
				return nil
			},
			StoreFunc: func(profile string, key string, value string) {
				is.Equal(profile, "profile")
				stored[key] = value
			},
		}

		var b bytes.Buffer
		cmd := Import(cfg)
		cmd.SetOut(&b)
		cmd.SetArgs([]string{"profile", `{"user.email": "work@example.com", "core.autocrlf": "input"}`})
		err := cmd.Execute()

		is.NoErr(err)
		is.Equal(trim(b.String()), "Successfully imported `profile` profile.")
		is.Equal(stored, map[string]string{
			"user.email":    "work@example.com",
			"core.autocrlf": "input",
		})
	})

	t.Run("array format", func(t *testing.T) {
		is := is.New(t)

		stored := make(map[string]string)
		cfg := &storageMock{
			SaveFunc: func(filename string) error {
				return nil
			},
			StoreFunc: func(profile string, key string, value string) {
				is.Equal(profile, "profile")
				stored[key] = value
			},
		}

		var b bytes.Buffer
		cmd := Import(cfg)
		cmd.SetOut(&b)
		cmd.SetArgs([]string{"profile", `[{"key": "user.email", "value": "work@example.com"}, {"key": "core.autocrlf", "value": "input"}]`})
		err := cmd.Execute()

		is.NoErr(err)
		is.Equal(trim(b.String()), "Successfully imported `profile` profile.")
		is.Equal(stored, map[string]string{
			"user.email":    "work@example.com",
			"core.autocrlf": "input",
		})
	})
}
