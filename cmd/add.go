package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// NewAdd returns `add` command
func NewAdd(c *Cmd) *cobra.Command {
	return &cobra.Command{
		Use:     "add [profile] [key] [value]",
		Aliases: []string{"set"},
		Short:   "Add an entry to a profile",
		Long:    "Adds an entry to a profile or update exists profile.",
		Example: "git-profile add my-profile user.email work@example.com\ngit-profile add my-profile user.name \"John Doe\"\ngit-profile add my-profile user.signingkey AAAAAAAA",
		Args:    cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			profile := args[0]
			key := args[1]
			value := args[2]

			c.storage.Store(profile, key, value)
			err := c.storage.Save(c.filename)
			if err != nil {
				cmd.PrintErr("Unable to store config file\n", err)
				os.Exit(1)
			}

			cmd.Printf("Successfully added `%s=%s` to `%s` profile.", key, value, profile)
		},
	}
}
