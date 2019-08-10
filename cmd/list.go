package cmd

import (
	"github.com/spf13/cobra"
)

// NewList returns `list` command
func NewList(c *Cmd) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "List of profiles",
		Long:    "Displays a list of available profiles.",
		Run: func(cmd *cobra.Command, args []string) {
			if c.storage == nil || len(c.storage.Profiles) == 0 {
				cmd.Println(`There are no available profiles.`)
				cmd.Println(`To add a profile, use the following examples:`)
				cmd.Println(`  git-profile add my-profile user.name "John Doe"`)
				cmd.Println(`  git-profile add my-profile user.email work@example.com`)
				return
			}

			cmd.Print("Available profiles:\n\n")
			for profile := range c.storage.Profiles {
				cmd.Printf("[%s]\n", profile)
				entries := c.storage.Profiles[profile]
				for _, entry := range entries {
					cmd.Printf("\t%s = %s\n", entry.Key, entry.Value)
				}
			}
		},
	}
}
