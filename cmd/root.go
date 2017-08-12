package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "git-profile",
	Short: "Allows you to switch between multiple user profiles in git repositories",
	Long: `Git Profile allows you to add and switch between multiple
user profiles in your git repositories.`,
}

// Execute executes the root command.
func Execute() {
	rootCmd.Execute()
}
