package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/kdisneur/changelog/pkg/git"
)

var remoteHTTPSExtractor = regexp.MustCompile("https://([^/]+)/(.*)")
var remoteGitExtractor = regexp.MustCompile("git@([^:]+):(.*)")

func FindRemoteFromURLs(remoteURLs []string) (*git.Remote, error) {
	if len(remoteURLs) == 0 {
		return nil, errors.New("no remote available")
	}

	remoteURL := remoteURLs[0]

	for _, newRemoteURL := range remoteURLs {
		if remoteURL != newRemoteURL {
			return nil, fmt.Errorf("found multiple remotes: %s, %s", remoteURL, newRemoteURL)
		}
	}

	return remoteFromURL(remoteURL)
}

func remoteFromURL(url string) (*git.Remote, error) {
	if strings.HasPrefix(url, "https://") {
		return remoteFromHTTPSURL(url)
	} else if strings.HasPrefix(url, "git@") {
		return remoteFromGitURL(url)
	}

	return nil, fmt.Errorf("unrecognized Git protocol for %s", url)
}

func remoteFromHTTPSURL(url string) (*git.Remote, error) {
	matches := remoteHTTPSExtractor.FindStringSubmatch(url)

	if len(matches) == 3 {
		return &git.Remote{Type: git.HTTPS, Host: matches[1], RepositoryName: matches[2]}, nil
	}

	return nil, fmt.Errorf("can't parse HTTPS remote: %s", url)
}

func remoteFromGitURL(url string) (*git.Remote, error) {
	matches := remoteGitExtractor.FindStringSubmatch(url)

	if len(matches) == 3 {
		return &git.Remote{Type: git.GIT, Host: matches[1], RepositoryName: matches[2]}, nil
	}

	return nil, fmt.Errorf("can't parse Git remote: %s", url)
}
