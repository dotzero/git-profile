package cmd

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"

	"github.com/dotzero/git-profile/config"
)

// Import returns `import` command
func Import(cfg *config.Config, filename *string) *cobra.Command {
	return &cobra.Command{
		Use:     "import [profile] [json-values]",
		Aliases: []string{"i"},
		Short:   "Import profile",
		Long:    "Import profile from json.",
		Args:    cobra.ExactArgs(2),
		Example: "cat my-profile.json | xargs -0 git-profile import my-profile",
		Run: func(cmd *cobra.Command, args []string) {
			profile := args[0]

			var entries []config.Entry

			err := json.Unmarshal([]byte(args[1]), &entries)
			if err != nil {
				cmd.PrintErr("Unable to decode profile values\n", err)
				os.Exit(1)
			}

			for _, entry := range entries {
				cfg.Store(profile, entry.Key, entry.Value)
			}

			err = cfg.Save(*filename)
			if err != nil {
				cmd.PrintErrln("Unable to save config file:", err)
				os.Exit(1)
			}

			cmd.Printf("Successfully import `%s` profile", profile)
		},
	}
}
