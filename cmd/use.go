package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/dotzero/git-profile/git"
)

// NewUse returns `use` command
func NewUse(c *Cmd) *cobra.Command {
	return &cobra.Command{
		Use:     "use [profile]",
		Aliases: []string{"u"},
		Short:   "Use a profile",
		Long:    "Applies selected profile entries to current git repository.",
		Example: "git-profile use my-profile",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if !git.IsRepository() {
				cmd.PrintErr("Current directory is not a root of git repository.")
				os.Exit(1)
			}

			profile := args[0]
			entries, ok := c.storage.Profiles[profile]
			if !ok {
				cmd.PrintErrf("There is no profile with `%s` name\n", profile)
				os.Exit(0)
			}

			err := git.Store(currentProfileKey, profile)
			if err != nil {
				cmd.PrintErr("Unable to interact with git to store current profile\n", err)
				os.Exit(1)
			}

			for _, entry := range entries {
				err := git.Store(entry.Key, entry.Value)
				if err != nil {
					cmd.PrintErr("Unable to interact with git to set profile entries\n", err)
					os.Exit(1)
				}
			}

			cmd.Printf("Successfully applied `%s` profile to current git repository.", profile)
		},
	}
}
