package cmd

import (
	"os"

	"github.com/dotzero/git-profile/git"
	"github.com/spf13/cobra"
)

// DefaultProfileName is a default profile name if not selected
const DefaultProfileName = `default`

var currentCmd = &cobra.Command{
	Use:     "current",
	Aliases: []string{"c"},
	Short:   "Show selected profile",
	Long:    "Show selected profile for current repository.",
	Run:     currentRun,
}

func init() {
	rootCmd.AddCommand(currentCmd)
}

func currentRun(cmd *cobra.Command, args []string) {
	if len(cfgStorage.Profiles) == 0 || !git.IsRepository() {
		os.Exit(1)
	}

	res, err := git.GetLocalConfig(`current-profile.name`)
	if len(res) == 0 || err != nil {
		cmd.Print(DefaultProfileName)
		os.Exit(0)
	}

	cmd.Printf("%s", res)
}
