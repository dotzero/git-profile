---
name: git-profile
description: Manage multiple Git identities with the git-profile CLI. Use when a user wants to create, inspect, apply, remove, import, export, or switch per-repository Git profiles (for example work and personal user.name, user.email, or signing key), or asks which git-profile identity is active.
---

# Git Profile

Use `git-profile` rather than manually editing `~/.gitprofile` or `.git/config`.

## Safety and scope

- Treat profile data as sensitive personal information. Do not echo values unnecessarily, especially signing keys.
- `add`, `del`, and `import` modify the profile store (by default `~/.gitprofile`).
- `use` and `unuse` modify the target repository's local Git configuration only; they never change global Git config.
- Before changing a repository, identify its target. Run `git rev-parse --show-toplevel` in that directory, or use `git-profile -C /path/to/repo ...` when the target is explicit.
- `git-profile list`, `use`, `unuse`, and `current` require a Git repository and at least one saved profile. Run `add`, `import`, or `export` from any directory.

## Check availability

Run `git-profile version` or `command -v git-profile` before using the tool. If it is unavailable, tell the user and offer one of the documented installation methods:

```sh
brew install dotzero/tap/git-profile
# or
go install github.com/dotzero/git-profile@latest
```

## Manage profile definitions

Create or update a profile non-interactively when the required values are known:

```sh
git-profile add work user.name "Jane Doe"
git-profile add work user.email jane.doe@company.example
git-profile add work user.signingkey ABCDEF0123456789
```

Use `git-profile add` with no arguments only when interactive input is appropriate. It offers the three common fields above; empty fields are skipped. The CLI also supports other valid Git config keys:

```sh
git-profile add work commit.gpgsign true
```

Inspect stored profiles in a repository:

```sh
git-profile list
```

Delete deliberately. `del profile` removes the whole stored profile; `del profile key` removes only that key:

```sh
git-profile del work user.signingkey
git-profile del work
```

For a non-default profile file, pass the global option explicitly:

```sh
git-profile --config /path/to/profiles.json add work user.email jane.doe@company.example
```

## Apply and verify a profile

Apply a named profile to the current repository:

```sh
git-profile use work
git-profile current
git config --local --get-regexp '^(user\.(name|email|signingkey)|current-profile\.name)$'
```

From elsewhere, put `-C` before the command:

```sh
git-profile -C /path/to/repo use work
git-profile -C /path/to/repo current
```

`use` records the selected profile in the local key `current-profile.name` and writes every entry stored in that profile. It does not remove local keys that were present in the previous profile but absent from the new one. When a clean switch matters, inspect the current profile and remove it first:

```sh
git-profile current
git-profile unuse old-profile
git-profile use new-profile
```

Do not run `unuse` blindly if the repository may have independent local values on the same keys: it unsets every key defined in the selected profile. To remove the applied profile, use:

```sh
git-profile unuse
```

## Transfer profiles

Export produces a JSON array suitable for import:

```sh
git-profile export work
git-profile import work '[{"key":"user.name","value":"Jane Doe"},{"key":"user.email","value":"jane.doe@company.example"}]'
```

Keep exported profile files private; they may contain personal email addresses and signing-key identifiers. Quote JSON as a single shell argument so its contents are passed unchanged.
