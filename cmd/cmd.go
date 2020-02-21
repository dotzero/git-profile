package cmd

import (
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"

	"github.com/dotzero/git-profile/config"
)

// Cmd is an CLI app
type Cmd struct {
	cobra.Command

	// Version is the version number or commit hash
	// These variables should be set by the linker when compiling
	Version string
	// CommitHash is the git hash of last commit
	CommitHash string
	// CompileDate is the date of build
	CompileDate string

	filename string
	storage  *config.Config
}

// New returns an pre-initialized CLI app
func New() *Cmd {
	root := cobra.Command{
		Use:   "git-profile",
		Short: "Allows you to switch between multiple user profiles in git repositories",
		Long:  "Git Profile allows to add and switch between multiple\nuser profiles in your git repositories.",
	}

	return &Cmd{
		Command:     root,
		Version:     "0.0.0-unknown",
		CommitHash:  "Unknown",
		CompileDate: "Unknown",
		storage:     config.New(),
	}
}

// Setup adds additional keys to CLI app
func (c *Cmd) Setup() {
	c.PersistentFlags().StringVarP(&c.filename, "config", "c", "~/.gitprofile", "config file")
	cobra.OnInitialize(func() {
		c.filename, _ = homedir.Expand(c.filename)
		err := c.storage.Load(c.filename)
		if err != nil {
			c.PrintErr("Unable to store config file\n", err)
			os.Exit(1)
		}
	})
}
