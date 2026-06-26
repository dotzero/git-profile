package cmd

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/matryer/is"

	"github.com/dotzero/git-profile/internal/config"
)

func TestUnuse(t *testing.T) {
	is := is.New(t)

	cfg := &storageMock{
		LenFunc: func() int {
			return 1
		},
		LookupFunc: func(name string) (config.Entry, bool) {
			return config.Entry{
				"user.email": "work@example.com",
			}, true
		},
	}

	var unset []string

	vcs := &vcsMock{
		IsRepositoryFunc: func() bool {
			return true
		},
		GetFunc: func(key string) (string, error) {
			if key == userEmailKey {
				return "work@example.com", nil
			}
			return "", nil
		},
		UnsetFunc: func(key string) error {
			unset = append(unset, key)
			return nil
		},
	}

	var b bytes.Buffer

	cmd := Unuse(cfg, vcs)

	cmd.SetOut(&b)
	cmd.SetArgs([]string{"profile"})
	err := cmd.Execute()

	is.NoErr(err)
	is.Equal(trim(b.String()), "Successfully removed `profile` profile from current git repository.")
	is.Equal(unset, []string{"user.email", currentProfileKey})
}

func TestUnuseProfileResolveArg(t *testing.T) {
	is := is.New(t)

	profile, err := unuseProfileResolve([]string{"work"}, &vcsMock{})

	is.NoErr(err)
	is.Equal(profile, "work")
}

func TestUnuseProfileResolveCurrent(t *testing.T) {
	is := is.New(t)

	vcs := &vcsMock{
		GetFunc: func(key string) (string, error) {
			is.Equal(key, currentProfileKey)
			return "home", nil
		},
	}

	profile, err := unuseProfileResolve(nil, vcs)

	is.NoErr(err)
	is.Equal(profile, "home")
}

func TestUnuseProfileResolveNoneApplied(t *testing.T) {
	is := is.New(t)

	vcs := &vcsMock{
		GetFunc: func(_ string) (string, error) {
			return "", fmt.Errorf("boom")
		},
	}

	_, err := unuseProfileResolve(nil, vcs)

	is.True(err != nil)
	is.Equal(err.Error(), "There is no profile applied to current git repository")
}
