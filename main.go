package main

import (
	"github.com/dotzero/git-profile/cmd"
)

var (
	// Version is the version number or commit hash
	// These variables should be set by the linker when compiling
	Version = "0.0.0-unknown"
	// CommitHash is the git hash of last commit
	CommitHash = "Unknown"
	// CompileDate is the date of build
	CompileDate = "Unknown"
)

func main() {
	c := cmd.New()
	c.Version = Version
	c.CommitHash = CommitHash
	c.CompileDate = CompileDate
	c.Execute()
}
