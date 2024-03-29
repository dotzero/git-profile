package cmd

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
)

// Export returns `export` command
func Export(cfg storage) *cobra.Command {
	return &cobra.Command{
		Use:     "export [profile]",
		Aliases: []string{"e"},
		Short:   "Export a profile",
		Long:    "Export a profile in JSON format.",
		Args:    cobra.ExactArgs(1),
		Example: "git-profile export my-profile",
		Run: func(cmd *cobra.Command, args []string) {
			profile := args[0]

			entries, ok := cfg.Lookup(profile)
			if !ok {
				cmd.PrintErrf("There is no profile with `%s` name\n", profile)
				os.Exit(0)
			}

			data, err := json.Marshal(entries)
			if err != nil {
				cmd.PrintErrln("Unable to encode profile values:", err)
				os.Exit(1)
			}

			cmd.Printf(string(data))
		},
	}
}
