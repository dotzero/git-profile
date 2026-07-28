package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/dotzero/git-profile/internal/config"
	"github.com/dotzero/git-profile/internal/ui"
)

// Migrate returns `migrate` command
func Migrate() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate ~/.gitprofile to the XDG config path",
		Long: multiline(
			"Copy ~/.gitprofile to $XDG_CONFIG_HOME/git-profile/config.json.",
			"The destination uses the map-based format.",
			"The legacy file is left in place. Once the XDG file exists, it takes priority.",
		),
		Example: multiline(
			`git-profile migrate`,
			`git-profile migrate --force`,
		),
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, _ []string) {
			src, dst, err := config.Migrate(force)
			if err != nil {
				ui.PrintErrln(cmd, ui.ErrorStyle, "%s", err)
				os.Exit(1)
			}

			ui.Println(cmd, ui.SuccessStyle, "Migrated `%s` to `%s`.", src, dst)
			ui.Println(cmd, ui.DefaultStyle, "Legacy file kept at `%s`.", src)
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "overwrite an existing XDG config file")

	return cmd
}
