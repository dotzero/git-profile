package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"

	"github.com/dotzero/git-profile/config"
	"github.com/dotzero/git-profile/git"
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
	config   *config.Config
	git      *git.Git
}

// New returns an app
func New() *Cmd {
	root := cobra.Command{
		Use:   "git-profile",
		Short: "Allows you to switch between multiple user profiles in git repositories",
		Long: multiline(
			"Git Profile allows you to add and switch between multiple",
			"user profiles in your git repositories.",
		),
	}

	return &Cmd{
		Command:     root,
		Version:     "0.0.0-unknown",
		CommitHash:  "Unknown",
		CompileDate: "Unknown",
		config:      config.New(),
		git:         git.New(),
	}
}

// Execute initialize the application and run it
func (c *Cmd) Execute() {
	c.init()

	err := c.Command.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (c *Cmd) init() {
	cobra.OnInitialize(func() {
		filename, err := homedir.Expand(c.filename)
		if err != nil {
			c.PrintErrln(err)
			os.Exit(1)
		}

		err = c.config.Load(filename)
		if err != nil {
			c.PrintErrln("Unable to load config:", err)
			os.Exit(1)
		}

		c.filename = filename
	})

	c.AddCommand(
		Add(c.config),
		Current(c.config, c.git),
		Del(c.config),
		List(c.config),
		Export(c.config),
		Import(c.config),
		Use(c.config, c.git),
		Version(c),
	)

	c.SetOutput(os.Stdout)
	c.SetErr(os.Stderr)

	c.PersistentFlags().StringVarP(&c.filename, "config", "c", "~/.gitprofile", "config file")
}

func multiline(lines ...string) string {
	return strings.Join(lines, "\n")
}
