package configuration_test

import (
	"os"
	"testing"

	homedir "github.com/mitchellh/go-homedir"

	"github.com/kdisneur/changelog/pkg/configuration"
)

func overrideHome(value string) func() {
	currentValue := os.Getenv("HOME")

	os.Setenv("HOME", value)

	return func() { os.Setenv("HOME", currentValue) }
}

func TestDefaultFolderPathWhenHomeIsFound(t *testing.T) {
	homedir.DisableCache = true
	defer func() { homedir.DisableCache = false }()

	callback := overrideHome("/my/home/path/")
	defer callback()

	expectedPath := "/my/home/path/.config"
	actualPath, err := configuration.DefaultFolderPath()

	if err != nil {
		t.Fatalf("It should not raise errors %#v", err)
	}

	if expectedPath != actualPath {
		t.Fatalf("Expected path: %s; Actual path: %s", expectedPath, actualPath)
	}
}

func TestDefaultFileName(t *testing.T) {
	expectedPath := "changelog"
	actualPath := configuration.DefaultFileName()

	if expectedPath != actualPath {
		t.Fatalf("Expected path: %s; Actual path: %s", expectedPath, actualPath)
	}
}

func TestDefaultFilePath(t *testing.T) {
	homedir.DisableCache = true
	defer func() { homedir.DisableCache = false }()

	callback := overrideHome("/my/home/path/")
	defer callback()

	expectedPath := "/my/home/path/.config/changelog"
	actualPath, err := configuration.DefaultFilePath()

	if err != nil {
		t.Fatalf("It should not raise errors %#v", err)
	}

	if expectedPath != actualPath {
		t.Fatalf("Expected path: %s; Actual path: %s", expectedPath, actualPath)
	}
}
