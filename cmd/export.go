package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// NewExport returns `export` command
func NewExport(c *Cmd) *cobra.Command {
	return &cobra.Command{
		Use:     "export [profile]",
		Aliases: []string{"e"},
		Short:   "Export profile",
		Long:    "Export profile as json.",
		Args:    cobra.ExactArgs(1),
		Example: "git-profile export my-profile",
		Run: func(cmd *cobra.Command, args []string) {
			profile := args[0]
			entries, ok := c.storage.Profiles[profile]
			if !ok {
				cmd.PrintErrf("There is no profile with `%s` name\n", profile)
				os.Exit(0)
			}

			bytes, err := json.Marshal(entries)

			if err != nil {
				cmd.PrintErrf("Unable encode values\n")
				os.Exit(1)
			}

			fmt.Println(string(bytes))
		},
	}
}
