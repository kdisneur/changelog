# Changelog

Generate a Changelog based on your Git history

Example of output from [Bamboo SMTP](https://github.com/fewlinesco/bamboo_smtp):

```
- Fix failing HexDoc redirection ([#79])
- Add attachment support ([#35])
- Apply rfc822_encode to headers(FROM, BCC, CC, TO) ([#75])
- Make the hostname (FQDN) configurable ([#74])
- Update Elixir, OTP and all deps to latest versions available ([#69])

[#35]: https://github.com/fewlinesco/bamboo_smtp/pull/35
[#69]: https://github.com/fewlinesco/bamboo_smtp/pull/69
[#74]: https://github.com/fewlinesco/bamboo_smtp/pull/74
[#75]: https://github.com/fewlinesco/bamboo_smtp/pull/75
[#79]: https://github.com/fewlinesco/bamboo_smtp/pull/79
```

## Usage

All the information should be available using the `-h` option.

```
$> changelog -h
Usage: changelog [-d|-h] <from-ref> [<repository-name [<end-ref>]]

Flags:
  -d           Debug mode
  -h           Display this usage text

Environment Variables:
  CHANGELOG_CONFIG_PATH:    Path where all the configuration is stored. Configuration contains
                            only one key for the moment: `github.token = <token>`
                            Currently set to "/Users/kdisneur/.config/changelog"
  CHANGELOG_DEFAULT_BRANCH: Default Git reference used to build the Changelog.

Arguments:
  from-ref:         Git references used as a starting point to build the Changelog
  repository-name:  By default, tries to get the remote url defined in Git config
  end-ref:          Reference branches to find merges from. If not set, it fallbacks
                    to the CHANGELOG_DEFAULT_BRANCH environment variables and, as a
                    last fallback, "develop"

Example:
  # Build a changelog with every merged Pull Request from v1.5.0
  changelog v1.5.0
  # Build a changelog with every merged Pull Request from v1.5.0 fetching Pull Requests from
  # johndoe/acme_corp
  changelog v1.5.0 johndoe/acme_corp
  # Build a changelog with every merged Pull Request from v1.5.0 fetching Pull Requests from
  # johndoe/acme_corp and using master as a reference branch
  changelog v1.5.0 johndoe/acme_corp master```
```

## Tests

```
bats test
```
