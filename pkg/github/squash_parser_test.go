package github_test

import (
	"strings"
	"testing"

	"github.com/kdisneur/changelog/pkg/github"
	"github.com/kdisneur/changelog/pkg/parser"
)

func TestSquashParserIsValidParser(t *testing.T) {
	squashParser := github.NewSquashParser()
	_, ok := squashParser.(parser.Parser)

	if !ok {
		t.Error("Squash parser doesn't implement Parser interface")
	}
}

func TestSquashParserFindID(t *testing.T) {
	testCases := []struct {
		Name          string
		CommitMessage string
		IsValid       bool
		ErrorMessage  string
		Expected      string
	}{
		{
			"Commit message with a valid format",
			"Add a nice feature (#1337)",
			true,
			"",
			"1337",
		},
		{
			"Commit message with an invalid format",
			"Add a nice feature #1337",
			false,
			"can't parse commit subject",
			"",
		},
		{
			"Commit message with no numbers",
			"Add a nice feature",
			false,
			"can't parse commit subject",
			"",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			parser := github.NewSquashParser()
			actual, err := parser.FindID(testCase.CommitMessage)

			if err != nil && testCase.IsValid {
				t.Fatalf("Expected no errors but got one. Received: %s", err.Error())
			}

			if err == nil && !testCase.IsValid {
				t.Fatalf("Expected an error but got none. Parsed ID '%s' from '%s'", actual, testCase.CommitMessage)
			}

			if err != nil && !strings.Contains(err.Error(), testCase.ErrorMessage) {
				t.Fatalf("Wrong error. Expected: %s. Received: %s", testCase.ErrorMessage, err.Error())
			}

			if actual != testCase.Expected {
				t.Fatalf("Wrong ID. Expected: %s. Received: %s", testCase.Expected, actual)
			}
		})
	}
}

func TestSquashParserKeepCommit(t *testing.T) {
	testCases := []struct {
		Name          string
		CommitMessage string
		IsValid       bool
	}{
		{
			"Commit message with a valid format",
			"Add a nice feature (#1337)",
			true,
		},
		{
			"Commit message with an invalid format",
			"Add a nice feature #1337",
			false,
		},
		{
			"Commit message with no numbers",
			"Add a nice feature",
			false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			parser := github.NewSquashParser()
			actual := parser.KeepCommit(testCase.CommitMessage)

			if testCase.IsValid && !actual {
				t.Fatalf("Expected commit to be kept but is not: '%s'", testCase.CommitMessage)
			}

			if !testCase.IsValid && actual {
				t.Fatalf("Expected commit to be rejected but is kept: '%s'", testCase.CommitMessage)
			}
		})
	}
}
