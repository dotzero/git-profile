# Git Profile switcher

[![build](https://github.com/dotzero/git-profile/actions/workflows/ci.yml/badge.svg)](https://github.com/dotzero/git-profile/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dotzero/git-profile)](https://goreportcard.com/report/github.com/dotzero/git-profile)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/dotzero/git-profile/blob/master/LICENSE)

Git Profile allows you to switch between multiple user profiles in git repositories

## Installation

If you are MacOS user, you can use [Homebrew](http://brew.sh/):

```bash
brew install dotzero/tap/git-profile
```

### Prebuilt binaries

Download the binary from the [releases](https://github.com/dotzero/git-profile/releases) page and place it under `$PATH` directory.

### Building from source

If your operating system does not have a binary release, but does run Go, you can build it from the source.

```bash
go get -u github.com/dotzero/git-profile
```

The binary will then be installed to `$GOPATH/bin` (or your `$GOBIN`).

## Usage

Adds an entry to a profile or updates an existing profile

```bash
git profile add home user.name dotzero
git profile add home user.email "me@dotzero.ru"
git profile add home user.signingkey AAAAAAAA
```

Displays a list of available profiles

```bash
git profile list
```

Applies the selected profile entries to the current git repository

```bash
git profile use home

# Under the hood it runs following commands:
# git config --local user.name dotzero
# git config --local user.email "me@dotzero.ru"
# git config --local user.signingkey AAAAAAAA
```

Export a profile in JSON format

```bash
git profile export home > home.json
```

Import profile from JSON format

```bash
cat home.json | xargs -0 git profile import home
```

## License

http://www.opensource.org/licenses/mit-license.php
