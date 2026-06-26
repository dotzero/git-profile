package cmd

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"

	"github.com/dotzero/git-profile/internal/config"
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
		Example: `git-profile import my-profile '[{"key":"user.email","value":"work@example.com"}]'`,
		Run: func(cmd *cobra.Command, args []string) {
			profile := args[0]
			filename, _ := cmd.Flags().GetString("config")

			var entries config.Entry

			err := json.Unmarshal([]byte(args[1]), &entries)
			if err != nil {
				ui.PrintErrln(cmd, ui.ErrorStyle, "Unable to decode profile values: %s", err)
				os.Exit(1)
			}

			for _, key := range []string{userNameKey, userEmailKey, userSigningKeyKey} {
				cfg.Store(profile, key, entries[key])
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
