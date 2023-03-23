package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// Use returns `use` command
func Use(cfg storage, v vcs) *cobra.Command {
	return &cobra.Command{
		Use:     "use [profile]",
		Aliases: []string{"u"},
		Short:   "Use a profile",
		Long:    "Applies the selected profile entries to the current git repository.",
		Example: "git-profile use my-profile",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if !v.IsRepository() {
				cmd.PrintErrln("The current working directory is not a git repository.")
				os.Exit(1)
			}

			profile := args[0]

			entries, ok := cfg.Lookup(profile)
			if !ok {
				cmd.PrintErrf("There is no profile with `%s` name\n", profile)
				os.Exit(0)
			}

			err := v.Set(currentProfileKey, profile)
			if err != nil {
				cmd.PrintErrln("Unable to interact with git to store current profile:", err)
				os.Exit(1)
			}

			for _, entry := range entries {
				err := v.Set(entry.Key, entry.Value)
				if err != nil {
					cmd.PrintErrln("Unable to interact with git to set profile entries:", err)
					os.Exit(1)
				}
			}

			cmd.Printf("Successfully applied `%s` profile to current git repository.\n", profile)
		},
	}
}
