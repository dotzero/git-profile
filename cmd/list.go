package cmd

import (
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List all profiles",
	Long:    `List all profiles.`,
	Run:     listRun,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func listRun(cmd *cobra.Command, args []string) {
	for title := range сfgStorage.Profiles {
		cmd.Printf("[%s]\n", title)
		entries, _ := сfgStorage.GetProfile(title)
		for _, entry := range entries {
			cmd.Printf("\t%s = \"%s\"\n", entry.Key, entry.Value)
		}
	}
}
