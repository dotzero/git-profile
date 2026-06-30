package cmd

import (
	"sort"

	"github.com/spf13/cobra"

	"github.com/dotzero/git-profile/internal/ui"
)

// List returns `list` command
func List(cfg storage, v vcs) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "List profiles",
		Long:    "Display the list of available profiles.",
		Example: "git-profile list",
		Run: func(cmd *cobra.Command, _ []string) {
			check(cmd, cfg, v)

			ui.Println(cmd, ui.InfoStyle, "Available profiles:")

			for _, name := range cfg.Names() {
				ui.Println(cmd, ui.SuccessStyle, "- %s:", name)

				profile, _ := cfg.Lookup(name)
				keys := make([]string, 0, len(profile))
				for key := range profile {
					keys = append(keys, key)
				}
				sort.Strings(keys)
				for _, key := range keys {
					ui.Println(cmd, ui.DefaultStyle, "  %s: %s", key, profile[key])
				}
			}
		},
	}
}
