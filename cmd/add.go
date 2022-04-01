package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/dotzero/git-profile/config"
)

// Add returns `add` command
func Add(cfg *config.Config, filename *string) *cobra.Command {
	return &cobra.Command{
		Use:     "add [profile] [key] [value]",
		Aliases: []string{"set"},
		Short:   "Add an entry to a profile",
		Long:    "Adds an entry to a profile or updates an existing profile.",
		Example: multiline(
			`git-profile add my-profile user.email work@example.com`,
			`git-profile add my-profile user.name "John Doe"`,
			`git-profile add my-profile user.signingkey AAAAAAAA`,
		),
		Args: cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			profile := args[0]
			key := args[1]
			value := args[2]

			cfg.Store(profile, key, value)

			err := cfg.Save(*filename)
			if err != nil {
				cmd.PrintErrln("Unable to save config file:", err)
				os.Exit(1)
			}

			cmd.Printf("Successfully stored `%s=%s` to `%s` profile.\n", key, value, profile)
		},
	}
}
