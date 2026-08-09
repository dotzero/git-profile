package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/dotzero/git-profile/internal/config"
	"github.com/dotzero/git-profile/internal/git"
)

// Cmd is a CLI app
type Cmd struct {
	cobra.Command

	// Version is the version number or commit hash
	// These variables should be set by the linker when compiling
	Version string
	// CommitHash is the git hash of last commit
	CommitHash string
	// CompileDate is the date of build
	CompileDate string

	filename  string
	directory string
	config    *config.Config
	git       *git.Git
}

// New returns an app
func New() *Cmd {
	root := cobra.Command{
		Use:   "git-profile",
		Short: "Manage multiple Git identities per repository",
		Long: multiline(
			"Git Profile helps you manage multiple Git identities",
			"and switch between them per repository.",
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
	c.PersistentPreRun = func(cmd *cobra.Command, _ []string) {
		if cmd.Name() == "migrate" {
			return
		}

		filename, err := resolveConfigPath(c.filename)
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
		c.git.SetDir(c.directory)
	}

	c.AddCommand(
		Add(c.config),
		Completion(&c.Command),
		Current(c.config, c.git),
		Del(c.config),
		List(c.config, c.git),
		Export(c.config),
		Import(c.config),
		Migrate(),
		Use(c.config, c.git),
		Unuse(c.config, c.git),
		Version(c),
	)

	c.SetOutput(os.Stdout)
	c.SetErr(os.Stderr)

	c.registerFlags()
}

func (c *Cmd) registerFlags() {
	c.PersistentFlags().StringVarP(
		&c.filename,
		"config",
		"c",
		"",
		"config file (default: $XDG_CONFIG_HOME/git-profile/config.json or ~/.gitprofile)",
	)
	c.PersistentFlags().StringVarP(&c.directory, "directory", "C", "", "run git commands in the given path")
}

func resolveConfigPath(filename string) (string, error) {
	if filename == "" {
		return config.DefaultPath()
	}

	return config.ExpandPath(filename)
}

func multiline(lines ...string) string {
	return strings.Join(lines, "\n")
}
