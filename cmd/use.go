package cmd

import (
	"os"
	"os/exec"

	"github.com/dotzero/git-profile/config"
	"github.com/spf13/cobra"
)

var useCmd = &cobra.Command{
	Use:     "use [profile]",
	Aliases: []string{"u"},
	Short:   "Use the profile",
	Long: `Use the profile.

Example:
  git-profile use work`,
	Run: useRun,
}

func init() {
	rootCmd.AddCommand(useCmd)
}

func isRepository() bool {
	s, err := os.Stat(".git")
	return err == nil && s.IsDir()
}

func gitSetConfig(key string, value string) error {
	if err := exec.Command("git", "config", "--local", key, value).Run(); err != nil {
		return err
	}
	return nil
}

func useRun(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		os.Exit(0)
	}

	if !isRepository() {
		cmd.Println("Current directory is not a root of git repository.")
		os.Exit(0)
	}

	c := config.NewConfig()
	c.Load(".gitprofile")

	entries, ok := c.GetProfile(args[0])
	if !ok {
		cmd.Printf("There is no profile with `%s` name", args[0])
		os.Exit(0)
	}

	cmd.Printf("Applying profile with name `%s`:\n\n", args[0])
	for _, entry := range entries {
		cmd.Printf("\t%s = \"%s\"\n", entry.Key, entry.Value)
		gitSetConfig(entry.Key, entry.Value)
	}
}
