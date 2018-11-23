package github_test

import (
	"strings"
	"testing"

	"github.com/kdisneur/changelog/pkg/github"
	"github.com/kdisneur/changelog/pkg/parser"
)

func TestMergeParserIsValidParser(t *testing.T) {
	mergeParser := github.NewMergeParser()
	_, ok := mergeParser.(parser.Parser)

	if !ok {
		t.Error("Merge parser doesn't implement Parser interface")
	}
}

func TestMergeParserFindID(t *testing.T) {
	testCases := []struct {
		Name          string
		CommitMessage string
		IsValid       bool
		ErrorMessage  string
		Expected      string
	}{
		{
			"Commit message with a valid format",
			"Merge pull request #1337 from kdisneur/changelog/#42-my_issue",
			true,
			"",
			"1337",
		},
		{
			"Commit message with an invalid format",
			"Merge pull request 1337 from kdisneur/changelog/#42-my_issue",
			false,
			"can't parse merge subject",
			"",
		},
		{
			"Commit message with no numbers",
			"Merge pull request from kdisneur/changelog/#42-my_issue",
			false,
			"can't parse merge subject",
			"",
		},
		{
			"Commit message with no merge message",
			"Add a feature #1137 from kdisneur/changelog/#42-my_issue",
			false,
			"can't parse merge subject",
			"",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			parser := github.NewMergeParser()
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

func TestMergeParserKeepCommit(t *testing.T) {
	testCases := []struct {
		Name          string
		CommitMessage string
		IsValid       bool
	}{
		{
			"Commit message with a valid format",
			"Merge pull request #1337 from kdisneur/changelog/#42-my_issue",
			true,
		},
		{
			"Commit message with an invalid format",
			"Merge pull request 1337 from kdisneur/changelog/#42-my_issue",
			false,
		},
		{
			"Commit message with no numbers",
			"Merge pull request from kdisneur/changelog/#42-my_issue",
			false,
		},
		{
			"Commit message with no merge message",
			"Add a feature #1137 from kdisneur/changelog/#42-my_issue",
			false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			parser := github.NewMergeParser()
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
