#! /usr/bin/env bash

_folder=${1};
_firstTagName=${2};

rm -rf ${_folder};
mkdir -p ${_folder};
cd ${_folder};

git init .;

git commit --no-gpg-sign --allow-empty --message "Initial Commit";
git tag ${_firstTagName};

git checkout -b "feature/first";
git commit --no-gpg-sign --allow-empty --message "Implement feature 1";
git checkout -;
git merge  --no-gpg-sign --no-ff "feature/first" --message "Merge pull request #42 from feature/first";
git branch -d "feature/first"

git checkout -b "feature/second";
git commit --no-gpg-sign --allow-empty --message "Implement feature 2";
git checkout -;
git merge  --no-gpg-sign --no-ff "feature/second" --message "Merge pull request #1337 from feature/second";
git branch -d "feature/second"
