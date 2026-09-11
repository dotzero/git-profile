package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

const (
	completionCommandName = "completion"
	bashShell             = "bash"
	zshShell              = "zsh"
	fishShell             = "fish"
	powershellShell       = "powershell"
)

// Completion returns `completion` command
func Completion(rootCmd *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:   completionCommandName + " [bash|zsh|fish|powershell]",
		Short: "Generate shell completion script",
		Long: `Generate shell completion script for Git Profile.

To load completions:

Bash:

  $ source <(git-profile completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ git-profile completion bash > /etc/bash_completion.d/git-profile
  # macOS:
  $ git-profile completion bash > $(brew --prefix)/etc/bash_completion.d/git-profile

Zsh:

  # Create completion directory and generate file:
  $ mkdir -p ~/.zsh/completions
  $ git-profile completion zsh > ~/.zsh/completions/_git-profile

  # Add to ~/.zshrc if not already present:
  $ echo 'fpath=(~/.zsh/completions $fpath)' >> ~/.zshrc
  $ echo 'autoload -U compinit && compinit' >> ~/.zshrc

  # Reload shell or run: source ~/.zshrc

Fish:

  $ git-profile completion fish | source

  # To load completions for each session, execute once:
  $ git-profile completion fish > ~/.config/fish/completions/git-profile.fish
`,
		ValidArgs:             []string{bashShell, zshShell, fishShell, powershellShell},
		Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		DisableFlagsInUseLine: true,
		Run: func(_ *cobra.Command, args []string) {
			switch args[0] {
			case bashShell:
				_ = rootCmd.GenBashCompletionV2(os.Stdout, true)
			case zshShell:
				_ = rootCmd.GenZshCompletion(os.Stdout)
			case fishShell:
				_ = rootCmd.GenFishCompletion(os.Stdout, true)
			case powershellShell:
				_ = rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
			}
		},
	}
}
