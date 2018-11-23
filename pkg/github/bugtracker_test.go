package github_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/kdisneur/changelog/pkg/bugtracker"
	"github.com/kdisneur/changelog/pkg/github"
	githubtest "github.com/kdisneur/changelog/pkg/testing/github"
)

func TestIsBugTracker(t *testing.T) {
	githubTracker := github.NewBugTracker("<api_token>", "kdisneur/changelog")
	_, ok := githubTracker.(bugtracker.BugTracker)

	if !ok {
		t.Errorf("GitHub tracker doesn't implement the Bugtracker interface")
	}
}

const ValidAPIToken string = "aaaa-bbbb-cccc-dddd"
const ValidRepositoryName string = "kdisneur/changelog"
const ValidPullRequestNumber string = "42"
const ValidSubject string = "A good feature description"

func TestBugTrackerFindIssue(t *testing.T) {
	validGithubPullRequestNumber, _ := strconv.Atoi(ValidPullRequestNumber)

	testCases := []struct {
		Name              string
		Token             string
		Repository        string
		PullRequestNumber string
		Mock              githubtest.GitHubMock
		IsValid           bool
		ErrorMessage      string
		Expected          *bugtracker.Issue
	}{
		{
			"When pull-request exists",
			ValidAPIToken,
			ValidRepositoryName,
			ValidPullRequestNumber,
			githubtest.NewMock(ValidAPIToken, ValidRepositoryName, 20181120, validGithubPullRequestNumber, ValidSubject),
			true,
			"",
			&bugtracker.Issue{
				ID:      ValidPullRequestNumber,
				Subject: ValidSubject,
				Link:    fmt.Sprintf("https://github.com/%s/pulls/%s", ValidRepositoryName, ValidPullRequestNumber),
			},
		},
		{
			"When pull-request doesn't exist",
			ValidAPIToken,
			ValidRepositoryName,
			"wrong-number",
			githubtest.NewMock(ValidAPIToken, ValidRepositoryName, 20181120, validGithubPullRequestNumber, ValidSubject),
			false,
			"can't fetch pull request",
			nil,
		},
		{
			"When API token is invalid",
			"wrong-token",
			ValidRepositoryName,
			ValidPullRequestNumber,
			githubtest.NewMock(ValidAPIToken, ValidRepositoryName, 20181120, validGithubPullRequestNumber, ValidSubject),
			false,
			"can't fetch pull request",
			nil,
		},
		{
			"When repository doesn't exist",
			ValidAPIToken,
			"wrong-repository",
			ValidPullRequestNumber,
			githubtest.NewMock(ValidAPIToken, ValidRepositoryName, 20181120, validGithubPullRequestNumber, ValidSubject),
			false,
			"can't fetch pull request",
			nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(testCase.Mock.Handler))
			defer server.Close()

			githubTracker := github.NewBugTrackerWithAPI(testCase.Token, server.URL, testCase.Repository)
			actualIssue, err := githubTracker.FindIssue(testCase.PullRequestNumber)

			if err != nil && testCase.IsValid {
				t.Fatalf("Epexted no errors but got one. Received: %s", err.Error())
			}

			if err == nil && !testCase.IsValid {
				t.Fatalf("Epexted an error but got none. Received: %+v", actualIssue)
			}

			if !testCase.IsValid && !strings.Contains(err.Error(), testCase.ErrorMessage) {
				t.Fatalf("Wrong error message. Expected: %s\nReceived: %s", testCase.ErrorMessage, err.Error())
			}

			if testCase.IsValid && !testCase.Expected.Equal(actualIssue) {
				t.Fatalf("Wrong issue. Expected: %+v\nReceived: %+v", testCase.Expected, actualIssue)
			}
		})
	}
}
