package cmd

import (
	"github.com/spf13/cobra"
)

// List returns `list` command
func List(cfg storage) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "List of profiles",
		Long:    "Displays a list of available profiles.",
		Run: func(cmd *cobra.Command, args []string) {
			if cfg.Len() == 0 {
				cmd.Println(`There are no available profiles.`)
				cmd.Println(`To add a new profile, use the following examples:`)
				cmd.Println(`  git-profile add my-profile user.name "John Doe"`)
				cmd.Println(`  git-profile add my-profile user.email work@example.com`)
				return
			}

			cmd.Println("Available profiles:")
			for _, name := range cfg.Names() {
				cmd.Printf("[%s]\n", name)

				profile, _ := cfg.Lookup(name)
				for _, entry := range profile {
					cmd.Printf("\t%s = %s\n", entry.Key, entry.Value)
				}
			}
		},
	}
}
