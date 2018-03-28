#! /usr/bin/env bash

_folder=${1};
_repository=${2};

rm -rf ${_folder};
mkdir -p ${_folder};
cd ${_folder};

git init .;

git remote add origin git@github.com:${_repository};
