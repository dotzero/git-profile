package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/dotzero/git-profile/internal/ui"
)

// Use returns `use` command
func Use(cfg storage, v vcs) *cobra.Command {
	return useCommand(cfg, v, ui.SelectProfile)
}

func useCommand(
	cfg storage,
	v vcs,
	selectProfile func([]string, io.Reader, io.Writer) (string, error),
) *cobra.Command {
	return &cobra.Command{
		Use:     "use [profile]",
		Aliases: []string{"u"},
		Short:   "Use a profile",
		Long:    "Applies the selected profile entries to the current git repository.",
		Example: multiline(
			`git-profile use`,
			`git-profile use my-profile`,
		),
		PreRun: func(cmd *cobra.Command, _ []string) {
			check(cmd, cfg, v)
		},
		Run: func(cmd *cobra.Command, args []string) {
			profile, err := profileResolve(args, cmd.InOrStdin(), cmd.OutOrStdout(), cfg, selectProfile)
			if err != nil {
				ui.PrintErrln(cmd, ui.ErrorStyle, "%s", err)
				os.Exit(1)
			}

			profileApply(cmd, cfg, v, profile)
		},
	}
}

func profileResolve(
	args []string,
	in io.Reader,
	out io.Writer,
	cfg storage,
	selectProfile func([]string, io.Reader, io.Writer) (string, error),
) (string, error) {
	if len(args) > 0 {
		return args[0], nil
	}

	profile, err := selectProfile(cfg.Names(), in, out)
	if err != nil {
		return "", fmt.Errorf("Unable to select a profile: %w", err)
	}

	return profile, nil
}

func profileApply(cmd *cobra.Command, cfg storage, v vcs, profile string) {
	entries, ok := cfg.Lookup(profile)
	if !ok {
		ui.PrintErrln(cmd, ui.ErrorStyle, "There is no profile with `%s` name", profile)
		os.Exit(0)
	}

	err := v.Set(currentProfileKey, profile)
	if err != nil {
		ui.PrintErrln(cmd, ui.ErrorStyle, "Unable to interact with git to store current profile: %s", err)
		os.Exit(1)
	}

	for key, value := range entries {
		err := v.Set(key, value)
		if err != nil {
			ui.PrintErrln(cmd, ui.ErrorStyle, "Unable to interact with git to store current profile: %s", err)
			os.Exit(1)
		}
	}

	ui.Println(cmd, ui.SuccessStyle, "Successfully applied `%s` profile to current git repository.", profile)
}
