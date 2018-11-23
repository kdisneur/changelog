package utils

import (
	"os/exec"
)

func ExecCommand(path string, gitArgs ...string) (string, error) {
	gitDefaultCommand := []string{"-C", path}
	gitCommand := append(gitDefaultCommand, gitArgs...)

	result, err := exec.Command("git", gitCommand...).Output()

	if err != nil {
		return "", err
	}

	return string(result), nil
}
