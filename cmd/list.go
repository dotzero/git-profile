package cmd

import (
	"github.com/spf13/cobra"

	"github.com/dotzero/git-profile/config"
)

// List returns `list` command
func List(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "List of profiles",
		Long:    "Displays a list of available profiles.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(cfg.Profiles) == 0 {
				cmd.Println(`There are no available profiles.`)
				cmd.Println(`To add a new profile, use the following examples:`)
				cmd.Println(`  git-profile add my-profile user.name "John Doe"`)
				cmd.Println(`  git-profile add my-profile user.email work@example.com`)
				return
			}

			cmd.Println("Available profiles:")
			for profile := range cfg.Profiles {
				cmd.Printf("[%s]\n", profile)
				for _, entry := range cfg.Profiles[profile] {
					cmd.Printf("\t%s = %s\n", entry.Key, entry.Value)
				}
			}
		},
	}
}
