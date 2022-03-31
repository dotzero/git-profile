package cmd

import (
	"encoding/json"
	"github.com/dotzero/git-profile/config"
	"github.com/spf13/cobra"
	"os"
)

// NewImport returns `import` command
func NewImport(c *Cmd) *cobra.Command {
	return &cobra.Command{
		Use:     "import [profile] [json-values]",
		Aliases: []string{"i"},
		Short:   "Import profile",
		Long:    "Import profile as json.",
		Args:    cobra.ExactArgs(2),
		Example: "cat ./my-profile.json | xargs -0 git-profile import my-profile",
		Run: func(cmd *cobra.Command, args []string) {
			profile := args[0]
			_, ok := c.storage.Profiles[profile]
			if ok {
				cmd.PrintErrf("There is already have profile with `%s` name\n", profile)
				os.Exit(1)
			}

			var data []config.Entry

			err := json.Unmarshal([]byte(args[1]), &data)

			if err != nil {
				cmd.PrintErrf("Unable decode values\n")
				os.Exit(1)
			}

			for _, entry := range data {
				c.storage.Store(profile, entry.Key, entry.Value)
			}

			err = c.storage.Save(c.filename)
			if err != nil {
				cmd.PrintErr("Unable to store config file\n", err)
				os.Exit(1)
			}

			cmd.Printf("Successfully import `%s` profile", profile)
		},
	}
}
