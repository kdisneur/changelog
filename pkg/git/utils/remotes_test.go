package utils_test

import (
	"strings"
	"testing"

	"github.com/kdisneur/changelog/pkg/git"
	"github.com/kdisneur/changelog/pkg/git/utils"
)

func TestRemoteURLs(t *testing.T) {
	testCases := []struct {
		Name         string
		URLs         []string
		IsValid      bool
		ErrorMessage string
		Remote       *git.Remote
	}{
		{
			"When has one HTTPS remote URL",
			[]string{"https://anurl.com/user/repo"},
			true,
			"",
			&git.Remote{Type: git.HTTPS, Host: "anurl.com", RepositoryName: "user/repo"},
		},
		{
			"When has one GIT remote URL",
			[]string{"git@anurl.com:user/repo"},
			true,
			"",
			&git.Remote{Type: git.GIT, Host: "anurl.com", RepositoryName: "user/repo"},
		},
		{
			"When has multiple different remote URLs",
			[]string{"git@anurl.com:user/repo", "https://anotherurl.com/user/repo"},
			false,
			"multiple remotes",
			nil,
		},
		{
			"When has no remote URLs",
			[]string{},
			false,
			"no remote available",
			nil,
		},
		{
			"When has unsupported remote URL scheme",
			[]string{"ftp://anurl.com/user/repo"},
			false,
			"unrecognized Git protocol",
			nil,
		},
		{
			"When has multiple times the same remote URL",
			[]string{"git@anurl.com:user/repo", "git@anurl.com:user/repo"},
			true,
			"",
			&git.Remote{Type: git.GIT, Host: "anurl.com", RepositoryName: "user/repo"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			remote, err := utils.FindRemoteFromURLs(testCase.URLs)
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
					t.Errorf("Wrong remote host: Expected %v, Received: %v", remote.Type, testCase.Remote.Type)
				}

				if remote.RepositoryName != testCase.Remote.RepositoryName {
					t.Errorf("Wrong remote repository name: Expected %v, Received: %v", remote.Type, testCase.Remote.Type)
				}
			}
		})
	}
}
