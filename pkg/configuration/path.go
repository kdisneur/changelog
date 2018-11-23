package configuration

import (
	"path"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

const DEFAULT_NAME = "changelog"

func DefaultFolderPath() (string, error) {
	folder, err := homedir.Expand("~/.config")
	if err != nil {
		return "", errors.Wrap(err, "Can't find home folder")
	}

	return folder, nil
}

func DefaultFileName() string {
	return DEFAULT_NAME
}

func DefaultFilePath() (string, error) {
	folder, err := DefaultFolderPath()
	if err != nil {
		return "", err
	}

	return path.Join(folder, DefaultFileName()), nil
}
