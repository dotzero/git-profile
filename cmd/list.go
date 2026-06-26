package cmd

import (
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
				for _, key := range []string{userNameKey, userEmailKey, userSigningKeyKey} {
					if profile[key] == "" {
						continue
					}
					ui.Println(cmd, ui.DefaultStyle, "  %s: %s", key, profile[key])
				}
			}
		},
	}
}
