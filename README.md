# Changelog

Build a changelog based on Git and a Tracker.

The command line takes two arguments:

1. A reference to a git object we want to build the changelog from
2. The name of the version

```bash
$ changelog v1.4.0 v1.5.0
## v1.5.0 - 2018-11-23

- Handle username and password required error ([#102])
- Add authentication option ([#89])

[#102]: https://github.com/fewlinesco/bamboo_smtp/pull/102
[#89]: https://github.com/fewlinesco/bamboo_smtp/pull/89
```

## Installation

```
brew install kdisneur/homebrew-formulae/changelog
```

## Configuration Option

### File

A `toml` file can be created at `~/.config/changelog.toml`

Here the definition:
```toml
[general]
mergeStrategy = "squash" # the default strategy to use when parsing a git history
                         # it can be either: squash or merge. By default: squash

baseBranch = "develop" # the main git branch you merge to. By default: `master`

[github]
token = "<api-key>" # a personal access-token to fetch pull-requests description.

[[repository]]
name = "kdisneur/changelog" # name of the repository. By default it extracts the
                            # information from the git remote

mergeStrategy = "squash" # the default strategy to use when parsing a git history
                         # it can be either: squash or merge. By default: squash
                         # it overrides the [general] section

baseBranch = "master" # the main git branch you merge to. It overrides the [general]
                      # section

[[repository]]
name = "fewlinesco/bamboo_smtp"
mergeStrategy = "merge"
baseBranch = "develop"
```

### Command Line

The command line have some options:

- `--branch`: name of the base branch. It overrides anything defined in the `file`
  section.
- `--change-dir` path to the local git repository if the command is run outside the
  repository root path.
- `--config` path to a configuration file if different from `~/.config/changelog`
- `--repository` name of the GitHub repository. By default, it tries to read from the
  git remote
- `--strategy` the default strategy to use when parsing a git history. It can be
  either: squash or merge and overrides anything defined in the `file` section.

## Development

### Installation

```bash
go mod download
go build
```

### Testing

```bash
go test ./...
```

### Deployment

> :warning: Make sure the last version has been build before running the command.
> And don't forget to bump the VERSION file.

```bash
./scripts/release.sh
```
