package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/dotzero/git-profile/internal/ui"
)

func check(cmd *cobra.Command, cfg storage, v vcs) {
	checkRepository(cmd, v)

	if cfg.Len() == 0 {
		ui.PrintErrln(cmd, ui.ErrorStyle, `There are no available profiles.`)
		ui.Println(cmd, ui.SuccessStyle, `To add a new profile, use the following command:`)
		ui.Println(cmd, ui.DefaultStyle, `$ git-profile add`)
		os.Exit(0)
	}
}

func checkRepository(cmd *cobra.Command, v vcs) {
	if !v.IsRepository() {
		ui.PrintErrln(cmd, ui.ErrorStyle, "The current working directory is not a valid git repository.")
		os.Exit(1)
	}
}
