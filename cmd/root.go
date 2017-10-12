package cmd

import (
	"log"
	"os"

	"../config"
	"github.com/hashicorp/logutils"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var (
	// Version is the version number or commit hash
	// These variables should be set by the linker when compiling
	Version = "0.0.0-unknown"
	// CommitHash is the git hash of last commit
	CommitHash = "Unknown"
	// CompileDate is the date of build
	CompileDate = "Unknown"

	cfgStorage *config.Config
	cfgFile    string
	isDebug    bool

	rootCmd = &cobra.Command{
		Use:   "git-profile",
		Short: "Allows you to switch between multiple user profiles in git repositories",
		Long: `Git Profile allows to add and switch between multiple
user profiles in your git repositories.`,
	}
)

func init() {
	cobra.OnInitialize(initLogs, initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "~/.gitprofile", "config file")
	rootCmd.PersistentFlags().BoolVarP(&isDebug, "debug", "d", false, "show debug log")
}

// Execute executes the root command.
func Execute() {
	rootCmd.Execute()
}

func initLogs() {
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel("WARN"),
		Writer:   os.Stderr,
	}

	if isDebug {
		filter.MinLevel = logutils.LogLevel("DEBUG")
	}

	log.SetOutput(filter)
}

func initConfig() {
	cfgFile, _ = homedir.Expand(cfgFile)
	cfgStorage = config.NewConfig()

	err := cfgStorage.Load(cfgFile)
	if err != nil {
		log.Println("[ERROR] Cannot load json from", cfgFile, err)
		os.Exit(1)
	}
}
