package cmd

import (
	"os"

	"git-profile/git"
	"github.com/spf13/cobra"
)

var useCmd = &cobra.Command{
	Use:     "use [profile]",
	Aliases: []string{"u"},
	Short:   "Use a profile",
	Long:    "Applies selected profile entries to current git repository.",
	Example: "  git-profile use my-profile",
	Args:    cobra.ExactArgs(1),
	Run:     useRun,
}

func init() {
	rootCmd.AddCommand(useCmd)
}

func useRun(cmd *cobra.Command, args []string) {
	if !git.IsRepository() {
		cmd.Println("Current directory is not a root of git repository.")
		os.Exit(1)
	}

	profile := args[0]

	entries, ok := cfgStorage.GetProfile(profile)
	if !ok {
		cmd.Printf("There is no profile with `%s` name", profile)
		os.Exit(0)
	}

	for _, entry := range entries {
		git.SetLocalConfig(entry.Key, entry.Value)
	}

	cmd.Printf("Successfully applied `%s` profile to current git repository.", profile)
}
