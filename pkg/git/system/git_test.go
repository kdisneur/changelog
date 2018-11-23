package system_test

import (
	"strings"
	"testing"
	"time"

	"github.com/kdisneur/changelog/pkg/git"
	"github.com/kdisneur/changelog/pkg/git/system"
	"github.com/kdisneur/changelog/pkg/testing/targz"
)

func setupFixture(name string) (git.Git, func(), error) {
	destination, cleanup, err := targz.Untar(name)

	if err != nil {
		return nil, cleanup, err
	}

	repository, err := system.NewRepository(destination)

	return repository, cleanup, err
}

func TestNewRepositoryWhenPathDoesNotExist(t *testing.T) {
	repository, err := system.NewRepository("/a/wrong/path")
	if err == nil {
		t.Fatalf("Expected an error but got none.Repository: %v", repository)
	}
}

func TestNewRepositoryWhenPathExists(t *testing.T) {
	_, cleanup, err := setupFixture("squash")
	defer cleanup()

	if err != nil {
		t.Errorf("Expected no errors but got one: %v", err)
	}
}

func TestFindRemote(t *testing.T) {
	testCases := []struct {
		Name         string
		FixtureName  string
		IsValid      bool
		ErrorMessage string
		Remote       *git.Remote
	}{
		{
			"When has one HTTPS remote URL",
			"onehttpsremote",
			true,
			"",
			&git.Remote{Type: git.HTTPS, Host: "github.com", RepositoryName: "kdisneur/changelog"},
		},
		{
			"When has one GIT remote URL",
			"onegitremote",
			true,
			"",
			&git.Remote{Type: git.GIT, Host: "github.com", RepositoryName: "kdisneur/changelog"},
		},
		{
			"When has multiple different remote URLs",
			"multipledifferentremotes",
			false,
			"multiple remotes",
			nil,
		},
		{
			"When has no remote URLs",
			"noremotes",
			false,
			"no remote available",
			nil,
		},
		{
			"When has unsupported remote URL scheme",
			"unsupportedremotescheme",
			false,
			"unrecognized Git protocol",
			nil,
		},
		{
			"When has multiple times the same remote URL",
			"multiplesamescheme",
			true,
			"",
			&git.Remote{Type: git.GIT, Host: "github.com", RepositoryName: "kdisneur/changelog"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			repository, cleanup, err := setupFixture(testCase.FixtureName)
			defer cleanup()

			remote, err := repository.FindRemote()

			if err != nil && testCase.IsValid {
				t.Fatalf("Expected no errors but go one: %s", err.Error())
			}

			if err == nil && !testCase.IsValid {
				t.Fatalf("Expected errors but go none: %v", remote)
			}

			if !testCase.IsValid {
				if !strings.Contains(err.Error(), testCase.ErrorMessage) {
					t.Errorf("Wrong error: Expected to contain '%s', got: '%s'", testCase.ErrorMessage, err.Error())
				}
			}

			if testCase.IsValid {
				if remote.Type != testCase.Remote.Type {
					t.Errorf("Wrong remote type: Expected %v, Received: %v", remote.Type, testCase.Remote.Type)
				}

				if remote.Host != testCase.Remote.Host {
					t.Errorf("Wrong remote host: Expected %v, Received: %v", remote.Host, testCase.Remote.Host)
				}

				if remote.RepositoryName != testCase.Remote.RepositoryName {
					t.Errorf("Wrong remote repository name: Expected %v, Received: %v", remote.RepositoryName, testCase.Remote.RepositoryName)
				}
			}
		})
	}
}

func TestLog(t *testing.T) {
	secondsCETOfUTC := int((1 * time.Hour).Seconds())
	centralEuropeTime := time.FixedZone("CET", secondsCETOfUTC)

	testCases := []struct {
		Name            string
		FixtureName     string
		From            git.Reference
		To              git.Reference
		IsValid         bool
		ErrorMessage    string
		ExpectedCommits []*git.Commit
	}{
		{
			"When `from` reference happened after the `to` reference",
			"squash",
			git.Reference("master"),
			git.Reference("v1.0.0"),
			true,
			"",
			[]*git.Commit{},
		},
		{
			"When `from` reference doesn't exist",
			"squash",
			git.Reference("inexistent"),
			git.Reference("v1.0.0"),
			false,
			"Can't generate git logs for 'inexistent..v1.0.0'",
			[]*git.Commit{},
		},
		{
			"When `to` reference doesn't exist",
			"squash",
			git.Reference("v1.0.0"),
			git.Reference("inexistent"),
			false,
			"Can't generate git logs for 'v1.0.0..inexistent'",
			[]*git.Commit{},
		},
		{
			"When `from` and `to` exists",
			"squash",
			git.Reference("v1.0.0"),
			git.Reference("master"),
			true,
			"",
			[]*git.Commit{
				{
					ID:          "4f28c412c51c44c94daa3fced544567c3f94dd7b",
					Author:      git.Person{Fullname: "John Doe", Email: "john.doe@gmail.com"},
					AuthoredAt:  time.Date(2018, time.November, 17, 20, 33, 25, 0, centralEuropeTime),
					Committer:   git.Person{Fullname: "Kevin Disneur", Email: "kevin@disneur.me"},
					CommittedAt: time.Date(2018, time.November, 17, 20, 35, 21, 0, centralEuropeTime),
					IsMerge:     false,
					Message:     "Adding feature 4 (#777)",
				},
				{
					ID:          "a2bc4fd34ba164ad0c1a264340ce37b0dbdaa6ef",
					Author:      git.Person{Fullname: "Kevin Disneur", Email: "kevin@disneur.me"},
					AuthoredAt:  time.Date(2018, time.November, 17, 6, 35, 51, 0, centralEuropeTime),
					Committer:   git.Person{Fullname: "Kevin Disneur", Email: "kevin@disneur.me"},
					CommittedAt: time.Date(2018, time.November, 17, 6, 35, 51, 0, centralEuropeTime),
					IsMerge:     false,
					Message:     "Adding feature 3 (#1337)",
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			repository, cleanup, err := setupFixture(testCase.FixtureName)
			defer cleanup()

			commits, err := repository.Log(testCase.From, testCase.To)

			if err != nil && testCase.IsValid {
				t.Fatalf("Expected no errors but go one: %s", err.Error())
			}

			if err == nil && !testCase.IsValid {
				t.Fatalf("Expected errors but go none: %v", commits)
			}

			if !testCase.IsValid {
				if !strings.Contains(err.Error(), testCase.ErrorMessage) {
					t.Fatalf("Wrong error: Expected to contain '%s', got: '%s'", testCase.ErrorMessage, err.Error())
				}
			}

			if testCase.IsValid {
				if len(commits) != len(testCase.ExpectedCommits) {
					t.Fatalf("Expected to find %d commits, found %d.\nExpected: %+v\n Received: %+v", len(commits), len(testCase.ExpectedCommits), commits, testCase.ExpectedCommits)
				}

				for index, expectedCommit := range testCase.ExpectedCommits {
					actualCommit := commits[index]

					if !expectedCommit.Equal(actualCommit) {
						t.Errorf("Wrong Commit.\nExpected: %+v\nReceived: %+v", expectedCommit, actualCommit)
					}
				}
			}
		})
	}
}
