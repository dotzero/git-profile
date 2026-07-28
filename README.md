# Git Profile

[![build](https://github.com/dotzero/git-profile/actions/workflows/ci.yml/badge.svg)](https://github.com/dotzero/git-profile/actions/workflows/ci.yml)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/dotzero/git-profile/blob/master/LICENSE)

Git Profile helps you manage multiple Git identities and switch between them per repository.

![](./demo/demo.gif)

## Installation

### Homebrew

```bash
brew install dotzero/tap/git-profile
```

### Prebuilt binaries

Download a binary from the [releases](https://github.com/dotzero/git-profile/releases) page and place it in a directory listed in `$PATH`.

### Build from source

```bash
go install github.com/dotzero/git-profile@latest
```

The binary will be installed to `$GOBIN` or `$GOPATH/bin`.

## Quick start

Create a profile in interactive mode:

```bash
git-profile add
```

Apply a profile in the current Git repository:

```bash
git-profile use
```

## Usage

### `add`

Add a key to a profile or update an existing one:

```bash
git-profile add work user.name "John Doe"
git-profile add work user.email work@example.com
git-profile add work user.signingkey AAAAAAAA
```

Run without arguments to start interactive mode:

```bash
git-profile add
```

Interactive mode asks for:
- profile name
- `user.name`
- `user.email`
- `user.signingkey`

For an existing profile, current values are prefilled.
For a new profile, empty fields are skipped.

### `list`

Show all available profiles:

```bash
git-profile list
```

### `del`

Delete a profile or a specific key from a profile:

```bash
git-profile del
git-profile del work
git-profile del work user.email
```

When called without arguments, `del` opens an interactive profile selector and deletes the selected profile.
Press `Esc` or `Ctrl+C` to cancel.

### `use`

Apply a profile to the current Git repository:

```bash
git-profile use work
```

Under the hood, this writes local Git config values, for example:

```bash
git config --local user.name "John Doe"
git config --local user.email work@example.com
git config --local user.signingkey AAAAAAAA
```

Run without arguments to select a profile interactively:

```bash
git-profile use
```

`use` must be executed inside a Git repository.

### `unuse`

Remove the applied profile entries from the current Git repository:

```bash
git-profile unuse work
```

Under the hood, this unsets the local Git config values, for example:

```bash
git config --local --unset user.name
git config --local --unset user.email
git config --local --unset user.signingkey
```

Run without arguments to remove the currently applied profile:

```bash
git-profile unuse
```

`unuse` must be executed inside a Git repository.

### `current`

Show the currently selected profile for the current repository:

```bash
git-profile current
```

### `export`

Export a profile as JSON:

```bash
git-profile export work
git-profile export work > work.json
```

Example output:

```json
[
  { "key": "user.name", "value": "John Doe" },
  { "key": "user.email", "value": "work@example.com" },
  { "key": "user.signingkey", "value": "AAAAAAAA" }
]
```

### `import`

Import a profile from JSON:

```bash
git-profile import work '[{"key":"user.name","value":"John Doe"},{"key":"user.email","value":"work@example.com"}]'
```

If you already have JSON in a file:

```bash
git-profile import work "$(cat work.json)"
```

## Config file

By default, Git Profile looks for profiles in this order:

1. `$XDG_CONFIG_HOME/git-profile/config.json` (or `~/.config/git-profile/config.json`)
2. `~/.gitprofile` (legacy path and array-based format)

New installs create the XDG config file in the map-based format.
Existing `~/.gitprofile` files stay in the legacy array-based format and are not rewritten in place.

You can override the path with:

```bash
git-profile --config /path/to/file add work user.email work@example.com
```

## Shell Completion

Shell completion is available for `git-profile`.

Note: use `git-profile ...`, not `git profile ...`, when generating or using completion.

### Bash

Requires [bash-completion](https://github.com/scop/bash-completion).

```bash
source <(git-profile completion bash)
```

To install permanently:

```bash
git-profile completion bash | sudo tee /etc/bash_completion.d/git-profile
```

On macOS:

```bash
git-profile completion bash > $(brew --prefix)/etc/bash_completion.d/git-profile
```

### Zsh

```bash
mkdir -p ~/.zsh/completions
git-profile completion zsh > ~/.zsh/completions/_git-profile
echo 'fpath=(~/.zsh/completions $fpath)' >> ~/.zshrc
echo 'autoload -U compinit && compinit' >> ~/.zshrc
source ~/.zshrc
```

### Fish

```bash
git-profile completion fish | source
git-profile completion fish > ~/.config/fish/completions/git-profile.fish
```

## License

[MIT](http://www.opensource.org/licenses/mit-license.php)
