package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var delCmd = &cobra.Command{
	Use:     "del [profile] [key]",
	Aliases: []string{"rm"},
	Short:   "Delete an entry or profile",
	Long: `Delete an entry from profile or an entire profile.
Provide a "key" argument to remove only that key from profile.

Example:
  git-profile del my-profile -> to delete an entire profile
  git-profile del my-profile user.name -> to delete only user.name`,
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

	if len(args) == 2 {
		cmd.Print("Delete an entry from profile:\n\n")
		сfgStorage.RemoveValue(profile, args[1])
	} else {
		cmd.Print("Delete an entire profile:\n\n")
		сfgStorage.RemoveProfile(profile)
	}

	сfgStorage.Save(cfgFile)
}
