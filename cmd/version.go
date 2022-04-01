package cmd

import (
	"github.com/spf13/cobra"
)

// Version returns `version` command
func Version(c *Cmd) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Git Profile",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("Git Profile")
			cmd.Println("Version:", c.Version)
			cmd.Println("Commit hash:", c.CommitHash)
			cmd.Println("Compiled on", c.CompileDate)
		},
	}
}
