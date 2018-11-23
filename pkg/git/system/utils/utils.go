package utils

import (
	"os"
	"path"
)

func IsGitRepository(repositoryPath string) bool {
	_, err := os.Stat(path.Join(repositoryPath, ".git"))

	return err == nil || !os.IsNotExist(err)
}
