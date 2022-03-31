package cmd

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
)

// NewExport returns `export` command
func NewExport(c *Cmd) *cobra.Command {
	return &cobra.Command{
		Use:     "export [profile]",
		Aliases: []string{"e"},
		Short:   "Export a profile",
		Long:    "Export a profile as json.",
		Args:    cobra.ExactArgs(1),
		Example: "git-profile export my-profile",
		Run: func(cmd *cobra.Command, args []string) {
			profile := args[0]

			entries, ok := c.storage.Profiles[profile]
			if !ok {
				cmd.PrintErrf("There is no profile with `%s` name\n", profile)
				os.Exit(0)
			}

			data, err := json.Marshal(entries)
			if err != nil {
				cmd.PrintErr("Unable to encode profile values\n", err)
				os.Exit(1)
			}

			cmd.Printf("%s", string(data))
		},
	}
}
