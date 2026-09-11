package cmd

import (
	"github.com/spf13/cobra"

	"github.com/dotzero/git-profile/internal/ui"
)

const versionCommandName = "version"

// Version returns `version` command
func Version(c *Cmd) *cobra.Command {
	return &cobra.Command{
		Use:     versionCommandName,
		Short:   "Print version information",
		Long:    "Print the version, commit hash, and build date.",
		Example: "git-profile version",
		Run: func(cmd *cobra.Command, _ []string) {
			ui.Println(cmd, ui.InfoStyle, "Git Profile")
			ui.Println(cmd, ui.DefaultStyle, "Version: %s", c.Version)
			ui.Println(cmd, ui.DefaultStyle, "Commit hash: %s", c.CommitHash)
			ui.Println(cmd, ui.DefaultStyle, "Compiled on %s", c.CompileDate)
		},
	}
}
