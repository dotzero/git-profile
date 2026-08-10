package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/dotzero/git-profile/internal/ui"
)

const (
	userNameKey       = "user.name"
	userEmailKey      = "user.email"
	userSigningKeyKey = "user.signingkey"
)

// Add returns `add` command
func Add(cfg storage) *cobra.Command {
	return addCommand(cfg, ui.PromptProfileName, ui.PromptProfileFields)
}

func addCommand(
	cfg storage,
	promptProfileName func(io.Reader, io.Writer) (string, error),
	editProfileFields func(ui.ProfileFormData, io.Reader, io.Writer) (ui.ProfileFormData, error),
) *cobra.Command {
	return &cobra.Command{
		Use:     "add [profile] [key] [value]",
		Aliases: []string{"set"},
		Short:   "Add or update profile entries",
		Long:    "Add a key to a profile or update an existing profile.",
		Example: multiline(
			`git-profile add`,
			`git-profile add my-profile user.email work@example.com`,
			`git-profile add my-profile user.name "John Doe"`,
			`git-profile add my-profile user.signingkey AAAAAAAA`,
		),
		Args: func(_ *cobra.Command, args []string) error {
			if len(args) == 0 || len(args) == 3 {
				return nil
			}

			return fmt.Errorf("accepts either no arguments for interactive mode or exactly 3 arguments")
		},
		Run: func(cmd *cobra.Command, args []string) {
			filename, _ := cmd.Flags().GetString("config")

			switch len(args) {
			case 3:
				err := profileUpdateEntry(cmd, cfg, filename, args[0], args[1], args[2])
				if err != nil {
					ui.PrintErrln(cmd, ui.ErrorStyle, "Unable to save config file: %s", err)
					os.Exit(1)
				}

				return
			case 0:
				err := profileUpdateEntries(cmd, cfg, filename, promptProfileName, editProfileFields)
				if err != nil {
					if ui.IsAborted(err) {
						ui.PrintErrln(cmd, ui.ErrorStyle, "Interactive add cancelled.")
						return
					}

					ui.PrintErrln(cmd, ui.ErrorStyle, "Unable to save config file: %s", err)
					os.Exit(1)
				}
			}
		},
	}
}

func profileUpdateEntry(
	cmd *cobra.Command,
	cfg storage,
	filename string,
	profile string,
	key string,
	value string,
) error {
	cfg.Store(profile, key, value)

	err := cfg.Save(filename)
	if err != nil {
		return err
	}

	ui.Println(cmd, ui.SuccessStyle, "Successfully stored `%s=%s` to `%s` profile.\n", key, value, profile)

	return nil
}

//nolint:funlen
func profileUpdateEntries(
	cmd *cobra.Command,
	cfg storage,
	filename string,
	promptProfileName func(io.Reader, io.Writer) (string, error),
	editProfileFields func(ui.ProfileFormData, io.Reader, io.Writer) (ui.ProfileFormData, error),
) error {
	profile, err := promptProfileName(cmd.InOrStdin(), cmd.OutOrStdout())
	if err != nil {
		return err
	}

	formData := ui.ProfileFormData{
		Profile: profile,
	}

	if entries, ok := cfg.Lookup(profile); ok {
		formData.UserName = entries[userNameKey]
		formData.UserEmail = entries[userEmailKey]
		formData.UserSigningKey = entries[userSigningKeyKey]
	}

	result, err := editProfileFields(formData, cmd.InOrStdin(), cmd.OutOrStdout())
	if err != nil {
		return err
	}

	currentValues, _ := cfg.Lookup(result.Profile)
	changed := 0

	values := map[string]string{
		userNameKey:       result.UserName,
		userEmailKey:      result.UserEmail,
		userSigningKeyKey: result.UserSigningKey,
	}

	for _, key := range []string{userNameKey, userEmailKey, userSigningKeyKey} {
		value := strings.TrimSpace(values[key])
		if value == "" {
			continue
		}

		if currentValues[key] == value {
			continue
		}

		cfg.Store(result.Profile, key, value)

		changed++
	}

	if changed == 0 {
		ui.Println(cmd, ui.SuccessStyle, "No profile changes to save.")
		return nil
	}

	err = cfg.Save(filename)
	if err != nil {
		return err
	}

	ui.Println(cmd, ui.SuccessStyle, "Successfully updated `%s` profile.\n", result.Profile)

	return nil
}
