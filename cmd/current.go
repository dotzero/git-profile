package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/dotzero/git-profile/git"
)

const (
	currentProfileKey  = "current-profile"
	defaultProfileName = "default"
)

// NewCurrent returns `current` command
func NewCurrent(c *Cmd) *cobra.Command {
	return &cobra.Command{
		Use:     "current",
		Aliases: []string{"c"},
		Short:   "Show selected profile",
		Long:    "Show selected profile for current repository.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(c.storage.Profiles) == 0 || !git.IsRepository() {
				os.Exit(1)
			}

			res, err := git.Lead(currentProfileKey)
			if len(res) == 0 || err != nil {
				cmd.Print(defaultProfileName)
				os.Exit(0)
			}

			cmd.Printf("%s", res)
		},
	}
}
