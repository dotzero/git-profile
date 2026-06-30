package cmd

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"

	"github.com/dotzero/git-profile/internal/config"
	"github.com/dotzero/git-profile/internal/oldconfig"
	"github.com/dotzero/git-profile/internal/ui"
)

// Import returns `import` command
func Import(cfg storage) *cobra.Command {
	return &cobra.Command{
		Use:     "import [profile] [json-values]",
		Aliases: []string{"i"},
		Short:   "Import a profile",
		Long:    "Import a profile from JSON.",
		Args:    cobra.ExactArgs(2),
		Example: `git-profile import my-profile '{"user.email":"work@example.com"}'`,
		Run: func(cmd *cobra.Command, args []string) {
			profile := args[0]
			filename, _ := cmd.Flags().GetString("config")

			var entries config.Entry

			err := json.Unmarshal([]byte(args[1]), &entries)
			if err != nil {
				// Try to unmarshal as old format (array)
				var oldEntries []oldconfig.OldEntry
				errOld := json.Unmarshal([]byte(args[1]), &oldEntries)
				if errOld != nil {
					ui.PrintErrln(cmd, ui.ErrorStyle, "Unable to decode profile values: %s", err)
					os.Exit(1)
				}
				entries = make(config.Entry)
				for _, entry := range oldEntries {
					entries[entry.Key] = entry.Value
				}
			}

			for key, val := range entries {
				cfg.Store(profile, key, val)
			}

			err = cfg.Save(filename)
			if err != nil {
				ui.PrintErrln(cmd, ui.ErrorStyle, "Unable to save config file: %s", err)
				os.Exit(1)
			}

			ui.Println(cmd, ui.SuccessStyle, "Successfully imported `%s` profile.", profile)
		},
	}
}
