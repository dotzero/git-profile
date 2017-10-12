package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var delCmd = &cobra.Command{
	Use:     "del [profile] [key]",
	Aliases: []string{"rm"},
	Short:   "Delete an entry or a profile",
	Long: `Delete an entry from a profile or an entire profile.
Provide a "key" argument to remove only one key from a profile.`,
	Example: `  git-profile del my-profile -> to delete an entire profile
  git-profile del my-profile user.name -> to delete only user.name`,
	Args: cobra.RangeArgs(1, 2),
	Run:  delRun,
}

func init() {
	rootCmd.AddCommand(delCmd)
}

func delRun(cmd *cobra.Command, args []string) {
	profile := args[0]

	if len(args) == 2 {
		key := args[1]
		if ok := cfgStorage.RemoveValue(profile, key); !ok {
			cmd.Printf("There is no profile with `%s` name", profile)
			os.Exit(0)
		}
		cfgStorage.Save(cfgFile)
		cmd.Printf("Successfully removed `%s` from `%s` profile.", key, profile)
		os.Exit(0)
	}

	cfgStorage.RemoveProfile(profile)
	cfgStorage.Save(cfgFile)
	cmd.Printf("Successfully removed `%s` profile.", profile)
}
