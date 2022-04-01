package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// Del returns `del` command
func Del(cfg storage) *cobra.Command {
	return &cobra.Command{
		Use:     "del [profile] [key]",
		Aliases: []string{"rm"},
		Short:   "Delete an entry or a profile",
		Long: multiline(
			`Delete an entry from a profile or an entire profile.`,
			`Enter the "key" argument to remove the exact key from the profile.`,
		),
		Example: multiline(
			`git-profile del my-profile (delete the entire profile)`,
			`git-profile del my-profile user.name (deleting only a certain key)`,
		),
		Args: cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			profile := args[0]
			filename, _ := cmd.Flags().GetString("config")

			if len(args) == 1 { // nolint
				if ok := cfg.DeleteProfile(profile); !ok {
					cmd.PrintErrln("There is no profile with given name")
					os.Exit(1)
				}

				err := cfg.Save(filename)
				if err != nil {
					cmd.PrintErrln("Unable to save config file:", err)
					os.Exit(1)
				}

				cmd.Printf("Successfully removed `%s` profile.\n", profile)
			} else {
				if ok := cfg.Delete(profile, args[1]); !ok {
					cmd.PrintErrln("There is no profile with given name")
					os.Exit(1)
				}

				err := cfg.Save(filename)
				if err != nil {
					cmd.PrintErrln("Unable to save config file:", err)
					os.Exit(1)
				}

				cmd.Printf("Successfully removed `%s` from `%s` profile.\n", args[1], profile)
			}
		},
	}
}
