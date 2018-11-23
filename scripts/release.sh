#! /usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

_applicationName="changelog";

_rootFolder=$(cd $(dirname ${0})/..; pwd);
_releaseFolder=${_rootFolder}/releases;
_configFile=${1:-~/.config/github};
_gitRepository="kdisneur/changelog";
_version=$(cat VERSION);
_currentSHA=$(git rev-parse HEAD);
_releasePath=${_releaseFolder}/${_applicationName};
_releaseArchiveName=${_applicationName}-${_version}.tgz;
_releaseArchivePath=${_releaseFolder}/${_releaseArchiveName};

if [ -z ${GITHUB_API_TOKEN:-} ]; then
  if [ ! -f ${_configFile} ]; then
    &>2 echo "You must have a GITHUB_API_TOKEN environment variable or ~/.config/github file available"
    exit 1
  else
    GITHUB_API_TOKEN=$(cat ${_configFile});
  fi
fi

echo "Build release..."
mkdir -p ${_releaseFolder};
(cd ${_rootFolder}; go build -o ${_releasePath});
(tar czf ${_releaseArchivePath} -C ${_releaseFolder} ${_applicationName} );

_shasum=$(shasum -a 256 ${_releaseArchivePath} | cut -f1 -d' ');

echo "Create release..."
_release=$(
  curl \
    --silent \
    --request POST \
    --header "Authorization: bearer ${GITHUB_API_TOKEN}" \
    --header "Content-Type: application/json" \
    --data '
      { "tag_name": "'${_version}'"
      , "target_commitish": "'${_currentSHA}'"
      , "body": "Release '${_version}'"
      , "name": "'${_version}'"
      , "draft": false
      , "prerelease": false
      }' \
    https://api.github.com/repos/${_gitRepository}/releases);

_uploadUrl=$(
  grep -oE "upload_url\":[ ]*\"([^\"]+)\"" <<< ${_release} \
    | grep -oE 'https:.*assets');

echo "Upload assets..."
_asset=$(
  curl \
    --silent \
    --request POST \
    --header "Authorization: bearer ${GITHUB_API_TOKEN}" \
    --header "Content-Type: application/x-gzip" \
    --data-binary @${_releaseArchivePath} \
    "${_uploadUrl}?name=${_releaseArchiveName}&label=Executable");

_downloadUrl=$(
  grep -oE "browser_download_url\":[ ]*\"([^\"]+)\"" <<< ${_asset} \
    | grep -oE 'https:.*tgz');

cat <<EOF
class Changelog < Formula
  desc "Generate changelog based on Git history"
  homepage "https://github.com/${_gitRepository}"
  url "${_downloadUrl}"
  sha256 "${_shasum}"
  def install
    bin.install "${_applicationName}"
  end
end
EOF
