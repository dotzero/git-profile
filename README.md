# Git Profile switcher

[![Build Status](https://travis-ci.org/dotzero/git-profile.svg?branch=master)](https://travis-ci.org/dotzero/git-profile)
[![Go Report Card](https://goreportcard.com/badge/github.com/dotzero/git-profile)](https://goreportcard.com/report/github.com/dotzero/git-profile)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/dotzero/git-profile/blob/master/LICENSE)

Git Profile allows to add and switch between multiple user profiles in your git repositories.

## Installation

If you are OSX user, you can use [Homebrew](http://brew.sh/):

```bash
brew install dotzero/tap/git-profile
```

### Prebuilt binaries

Download the binary from the [releases](https://github.com/dotzero/git-profile/releases) page and place it in `$PATH` directory.

### Building from source

If your operating system does not have a binary release, but does run Go, you can build from source.

Make sure that you have Go version 1.7 or greater and that your `GOPATH` env variable is set (I recommand setting it to `~/go` if you don't have one).

```bash
go get -u github.com/dotzero/git-profile
```

The binary will then be installed to `$GOPATH/bin` (or your `$GOBIN`).

## Usage

Add an entry to a profile

```bash
git profile add home user.name dotzero
git profile add home user.email "mail@dotzero.ru"
git profile add home user.signingkey AAAAAAAA
```

List of available profiles

```bash
git profile list
```

Apply the profile to current git repository

```bash
git profile use home

# Under the hood it runs following commands:
# git config --local user.name dotzero
# git config --local user.email "mail@dotzero.ru"
# git config --local user.signingkey AAAAAAAA
```

Export profile to file

```bash
git profile export home > home.json
```

Import profile from file

```bash
cat ./my-profile.json | xargs -0 git-profile import my-profile
```

## License

http://www.opensource.org/licenses/mit-license.php
