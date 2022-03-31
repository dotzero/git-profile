package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// NewDel returns `del` command
func NewDel(c *Cmd) *cobra.Command {
	return &cobra.Command{
		Use:     "del [profile] [key]",
		Aliases: []string{"rm"},
		Short:   "Delete an entry or a profile",
		Long:    "Delete an entry from a profile or an entire profile.\nProvide a \"key\" argument to remove only one key from a profile.",
		Example: "git-profile del my-profile -> to delete an entire profile\ngit-profile del my-profile user.name -> to delete only user.name",
		Args:    cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			profile := args[0]

			if len(args) == 2 { // nolint
				key := args[1]
				ok := c.storage.Delete(profile, key)
				if !ok {
					cmd.PrintErrf("There is no profile with `%s` name\n", profile)
					os.Exit(1)
				}

				err := c.storage.Save(c.filename)
				if err != nil {
					cmd.PrintErr("Unable to store config file", err)
					os.Exit(1)
				}

				cmd.Printf("Successfully removed `%s` from `%s` profile.\n", key, profile)
				os.Exit(0)
			}

			delete(c.storage.Profiles, profile)

			err := c.storage.Save(c.filename)
			if err != nil {
				cmd.PrintErr("Unable to store config file\n", err)
				os.Exit(1)
			}

			cmd.Printf("Successfully removed `%s` profile.", profile)
		},
	}
}
