package cmd

import (
	"log"
	"os"

	"github.com/dotzero/git-profile/config"
	"github.com/hashicorp/logutils"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var (
	cfgFile    string
	сfgStorage *config.Config
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
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.gitprofile)")
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
	сfgPath, err := homedir.Expand("~/.gitprofile")
	if err != nil {
		log.Println("[ERROR] Cannot obtain ~/.gitprofile")
		os.Exit(1)
	}

	сfgStorage = config.NewConfig()
	err = сfgStorage.Load(сfgPath)
	if err != nil {
		log.Println("[ERROR] Cannot load json from ~/.gitprofile", err)
		os.Exit(1)
	}
}
