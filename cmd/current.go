package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/dotzero/git-profile/git"
)

var currentCmd = &cobra.Command{
	Use:     "current",
	Aliases: []string{"c"},
	Short:   "Show current profile",
	Long:    "Show current profile for repo",
	Run:     currentRun,
}

func init() {
	rootCmd.AddCommand(currentCmd)
}

func currentRun(cmd *cobra.Command, args []string) {
	if len(cfgStorage.Profiles) == 0 {
		cmd.Print(`profiles not setted`)
		os.Exit(0)
	}

	if (!git.IsRepository()) {
		//cmd.Print(`this is not a repository`)
		os.Exit(1)
	}

	res, err := git.GetLocalConfig(`current-profile.name`)

	if (len(res) == 0 || err != nil){
		cmd.Print(`default`)
		os.Exit(0)
	}

	cmd.Printf("%s", res)
	os.Exit(0)
}
