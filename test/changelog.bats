#! /usr/bin/env bats

setup() {
  export PATH=${BATS_TEST_DIRNAME}/stubs:${PATH};
  load ${BATS_TEST_DIRNAME}/../src/changelog;
}

@test "_build_changelog" {
  local pullRequests=("12 A small issue" "1337 Another issue solved" "42 Description of an issue")

	run _build_changelog "johndoe/acme_corp" "${pullRequests[@]}"

	[ ${status} -eq 0 ]

  echo $output

  [ ${lines[0]} = "- A small issue ([#12])" ]
  [ ${lines[1]} = "- Another issue solved ([#1337])" ]
  [ ${lines[2]} = "- Description of an issue ([#42])" ]
  [ ${lines[3]} = "[#12]: https://github.com/johndoe/acme_corp/pull/12" ]
  [ ${lines[4]} = "[#42]: https://github.com/johndoe/acme_corp/pull/42" ]
  [ ${lines[5]} = "[#1337]: https://github.com/johndoe/acme_corp/pull/1337" ]
}

@test "_ensureGitHubTokenIsAvailable: when token is not" {
  run _ensureGitHubTokenIsAvailable ""

  [ ${status} -eq 1 ]

  [[ ${lines[0]} =~ "GitHub token is not available" \
     && ${lines[1]} =~ "Usage" ]]
}

@test "_ensureGitHubTokenIsAvailable: when token is" {
  run _ensureGitHubTokenIsAvailable "secret-token"

  [ ${status} -eq 0 ]
}

@test "_error: when usage not set" {
  run _error "A nice message"

  [ ${status} -eq 1 ]
  [ ${output} = "A nice message" ]
}

@test "_error: when usage set to false" {
  run _error "false" "A nice message"

  [ ${status} -eq 1 ]
  [ ${output} = "A nice message" ]
}

@test "_error: when usage set to true" {
  run _error "true" "A nice message"

  [ ${status} -eq 1 ]
  [ ${lines[0]} = "A nice message" ]
  [[ ${lines[1]} =~ "Usage:" ]]
}

@test "_errorWithUsage: display the usage" {
  run _errorWithUsage "A nice message"

  [ ${status} -eq 1 ]
  [ ${lines[0]} = "A nice message" ]
  [[ ${lines[1]} =~ "Usage:" ]]
}

@test "_filterOnPullRequests: should work" {
  local folder=${BATS_TEST_DIRNAME}/fixtures/repo_with_logs

  ${BATS_TEST_DIRNAME}/helpers/setup-repo-with-logs.sh ${folder} v1.0.0 &>/dev/null

  cd ${folder};

  run _filterOnPullRequests "v1.0.0" "master"

  [ ${status} -eq 0 ]

  [ ${lines[0]} = "1337" ]
  [ ${lines[1]} = "42" ]
}

@test "_fetchPullRequestData: should work" {
  run _fetchPullRequestData "aaa-bbb-ccc-ddd" "johndoe/acme_corp" 42 1337

  [ ${status} -eq 0 ]

  [ ${lines[0]} = "42 Description of an issue" ]
  [ ${lines[1]} = "1337 Another issue solved" ]
}

@test "_fetchPullRequestTitleFromGitHub: should work" {
  run _fetchPullRequestTitleFromGitHub "aaa-bbb-ccc-ddd" "johndoe/acme_corp" 42

  [ ${status} -eq 0 ]
  [ "${output}" = "Description of an issue" ]
}

@test "_findGitHubRepositoryName: when repository is already given" {
  run _findGitHubRepositoryName "johndoe/given_repo"

  [ ${status} -eq 0 ]
  [ "${output}" = "johndoe/given_repo" ]
}

@test "_findGitHubRepositoryName: when repository is not given but one remote is set in config" {
  local folder=${BATS_TEST_DIRNAME}/fixtures/repo_with_one_remote

  ${BATS_TEST_DIRNAME}/helpers/setup-mono-remote.sh ${folder} "johndoe/mono_remote" &>/dev/null

  cd ${folder};

  run _findGitHubRepositoryName ""

  [ ${status} -eq 0 ]
  [ "${output}" = "johndoe/mono_remote" ]
}

@test "_findGitHubRepositoryName: when repository is not given but several remotes are set in config" {
  local folder=${BATS_TEST_DIRNAME}/fixtures/repo_with_several_remotes

  ${BATS_TEST_DIRNAME}/helpers/setup-dual-remote.sh ${folder} "johndoe/first_remote" "johndoe/second_remote" &>/dev/null

  cd ${folder};

  run _findGitHubRepositoryName ""

  [ ${status} -eq 1 ]
  [[ "${lines[0]}" = "We can't automatically find a remote URL: 2 found." \
    && "${lines[1]}" = "You will have to set the repository manually." \
    && "${lines[2]}" =~ "Usage" ]]
}

@test "_readConfig: when file doesn't exist" {
  run _readConfig ${BATS_TEST_DIRNAME}/fixtures/non-existing-file "an.awesome.key"

  [ ${status} -eq 1 ]
  [[ "${lines[0]}" = "Config file not found." \
    && "${lines[1]}" =~ "Usage: " ]]
}

@test "_readConfig: when key doesn't exist but default value is provided" {
  run _readConfig ${BATS_TEST_DIRNAME}/fixtures/dummy-config "key.notFound" "I'm default"

  [ ${status} -eq 0 ]
  [ "${output}" = "I'm default" ]
}

@test "_readConfig: when key exists with spaces around the equal sign and default value is provided" {
  run _readConfig ${BATS_TEST_DIRNAME}/fixtures/dummy-config "key.withSpaces" "I'm default"

  [ ${status} -eq 0 ]
  [ "${output}" = "a nice value with spaces between = sign" ]
}

@test "_readConfig: when key exists with spaces around the equal sign" {
  run _readConfig ${BATS_TEST_DIRNAME}/fixtures/dummy-config "key.withSpaces"

  [ ${status} -eq 0 ]
  [ "${output}" = "a nice value with spaces between = sign" ]
}

@test "_readConfig: when key exists with no spaces around the equal sign and default value is provided" {
  run _readConfig ${BATS_TEST_DIRNAME}/fixtures/dummy-config "key.withNoSpaces" "I'm default"

  [ ${status} -eq 0 ]
  [ "${output}" = "it works even if stuck together" ]
}

@test "_readConfig: when key exists with no spaces around the equal sign" {
  run _readConfig ${BATS_TEST_DIRNAME}/fixtures/dummy-config "key.withNoSpaces"

  [ ${status} -eq 0 ]
  [ "${output}" = "it works even if stuck together" ]
}
