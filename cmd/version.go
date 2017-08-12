package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Git Profile",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Git Profile v1.0")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
