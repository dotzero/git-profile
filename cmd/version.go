package cmd

import (
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Git Profile",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("Git Profile")
		cmd.Println("Version:", Version)
		cmd.Println("Commit hash:", CommitHash)
		cmd.Println("Compiled on", CompileDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
