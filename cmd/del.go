package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var delCmd = &cobra.Command{
	Use:     "del [profile]",
	Aliases: []string{"d"},
	Short:   "Delete an entry from profile",
	Long: `Delete an entry from profile.

Example:
  git-profile del work`,
	Run: delRun,
}

func init() {
	rootCmd.AddCommand(delCmd)
}

func delRun(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		os.Exit(1)
	}

	profile := args[0]

	сfgStorage.RemoveProfile(profile)
	сfgStorage.Save(cfgFile)
}
