package cmd

import (
	"log"
	"os"

	"git-profile/config"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:     "add [profile] [key] [value]",
	Aliases: []string{"set"},
	Short:   "Add an entry to a profile",
	Long:    "Adds an entry to a profile or update exists profile.",
	Example: `  git-profile add my-profile user.email work@example.com
  git-profile add my-profile user.name "John Doe"
  git-profile add my-profile user.signingkey AAAAAAAA`,
	Args: cobra.ExactArgs(3),
	Run:  addRun,
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func addRun(cmd *cobra.Command, args []string) {
	profile := args[0]
	key := args[1]
	value := args[2]

	cfgStorage.SetValue(profile, config.Entry{Key: key, Value: value})
	err := cfgStorage.Save(cfgFile)
	if err != nil {
		log.Println("[ERROR] Cannot save json to", cfgFile, err)
		os.Exit(1)
	}

	cmd.Printf("Successfully added `%s=%s` to `%s` profile.", key, value, profile)
}
