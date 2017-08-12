package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/dotzero/git-profile/config"
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:     "set [profile] [key=value]",
	Aliases: []string{"s"},
	Short:   "Set an entry in profile",
	Long: `Sets a new entry to profile.

Example:
  git-profile set work user.email=work@example.com
  git-profile set work user.name="John Doe"`,
	Run: setRun,
}

func init() {
	rootCmd.AddCommand(setCmd)
}

func setRun(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		cmd.Usage()
		os.Exit(1)
	}

	profile := args[0]
	kv := strings.Split(args[1], "=")
	key, value := kv[0], kv[1]

	fmt.Println(profile, key, value)

	c := config.NewConfig()
	c.Load(".gitprofile")
	c.SetValue(profile, config.Entry{Key: key, Value: value})
	c.Save(".gitprofile")
}
