package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/dotzero/git-profile/internal/ui"
)

// Unuse returns `unuse` command
func Unuse(cfg storage, v vcs) *cobra.Command {
	return unuseCommand(cfg, v)
}

func unuseCommand(cfg storage, v vcs) *cobra.Command {
	return &cobra.Command{
		Use:     "unuse [profile]",
		Aliases: []string{"uu"},
		Short:   "Unuse a profile",
		Long:    "Removes the applied profile entries from the current git repository.",
		Example: multiline(
			`git-profile unuse`,
			`git-profile unuse my-profile`,
		),
		PreRun: func(cmd *cobra.Command, _ []string) {
			check(cmd, cfg, v)
		},
		Run: func(cmd *cobra.Command, args []string) {
			profile, err := unuseProfileResolve(args, v)
			if err != nil {
				ui.PrintErrln(cmd, ui.ErrorStyle, "%s", err)
				os.Exit(0)
			}

			profileUnapply(cmd, cfg, v, profile)
		},
	}
}

func unuseProfileResolve(args []string, v vcs) (string, error) {
	if len(args) > 0 {
		return args[0], nil
	}

	profile, err := v.Get(currentProfileKey)
	if len(profile) == 0 || err != nil {
		return "", fmt.Errorf("There is no profile applied to current git repository")
	}

	return profile, nil
}

func profileUnapply(cmd *cobra.Command, cfg storage, v vcs, profile string) {
	entries, ok := cfg.Lookup(profile)
	if !ok {
		ui.PrintErrln(cmd, ui.ErrorStyle, "There is no profile with `%s` name", profile)
		os.Exit(0)
	}

	for key := range entries {
		err := v.Unset(key)
		if err != nil {
			ui.PrintErrln(cmd, ui.ErrorStyle, "Unable to interact with git to remove current profile: %s", err)
			os.Exit(1)
		}
	}

	err := v.Unset(currentProfileKey)
	if err != nil {
		ui.PrintErrln(cmd, ui.ErrorStyle, "Unable to interact with git to remove current profile: %s", err)
		os.Exit(1)
	}

	ui.Println(cmd, ui.SuccessStyle, "Successfully removed `%s` profile from current git repository.", profile)
}
