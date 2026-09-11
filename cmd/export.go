package cmd

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"

	"github.com/dotzero/git-profile/internal/ui"
)

// Export returns `export` command
func Export(cfg storage) *cobra.Command {
	return &cobra.Command{
		Use:     "export [profile]",
		Aliases: []string{"e"},
		Short:   "Export a profile",
		Long:    "Export a profile as JSON.",
		Args:    cobra.ExactArgs(1),
		Example: "git-profile export my-profile",
		Run: func(cmd *cobra.Command, args []string) {
			profile := args[0]

			entries, ok := cfg.Lookup(profile)
			if !ok {
				ui.PrintErrln(cmd, ui.ErrorStyle, "There is no profile with `%s` name", profile)
				os.Exit(1)
			}

			data, err := json.Marshal(entries)
			if err != nil {
				ui.PrintErrln(cmd, ui.ErrorStyle, "Unable to encode profile values: %s", err)
				os.Exit(1)
			}

			cmd.Println(string(data))
		},
	}
}
