package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List of profiles",
	Long:    "Displays a list of available profiles.",
	Run:     listRun,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func listRun(cmd *cobra.Command, args []string) {
	if len(сfgStorage.Profiles) == 0 {
		cmd.Print(`There are no available profiles.
To add a profile, use the following examples:
  git-profile add my-profile user.name "John Doe"
  git-profile add my-profile user.email work@example.com`)
		os.Exit(0)
	}

	cmd.Print("Available profiles:\n\n")
	for title := range сfgStorage.Profiles {
		cmd.Printf("[%s]\n", title)
		entries, _ := сfgStorage.GetProfile(title)
		for _, entry := range entries {
			cmd.Printf("\t%s = %s\n", entry.Key, entry.Value)
		}
	}
}
