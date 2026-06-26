package cmd

//go:generate go tool github.com/matryer/moq -skip-ensure -out mock.go . storage vcs

import (
	"github.com/dotzero/git-profile/internal/config"
)

type storage interface {
	Len() int
	Lookup(name string) (config.Entry, bool)
	Names() []string
	Delete(profile string, value string) bool
	DeleteProfile(profile string) bool
	Store(profile string, key string, value string)
	Save(filename string) error
	Load(filename string) (err error)
}

type vcs interface {
	IsRepository() bool
	Get(key string) (string, error)
	Set(key string, value string) error
	Unset(key string) error
}
