package github

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/kdisneur/changelog/pkg/bugtracker"
	"github.com/pkg/errors"
)

func NewBugTracker(token string, repository string) bugtracker.BugTracker {
	return NewBugTrackerWithAPI(token, "https://api.github.com", repository)
}

func NewBugTrackerWithAPI(token string, apiURL string, repository string) bugtracker.BugTracker {
	return GitHub{token, apiURL, repository}
}

func (g GitHub) Equal(other bugtracker.BugTracker) bool {
	tracker, hasGoodType := other.(GitHub)
	if !hasGoodType {
		return false
	}

	return g.Token == tracker.Token && g.API_URL == tracker.API_URL && g.Repository == tracker.Repository
}

func (g GitHub) FindIssue(id string) (*bugtracker.Issue, error) {
	client := &http.Client{}

	request, err := http.NewRequest("GET", githubPullRequestPath(g, id), nil)
	if err != nil {
		return nil, errors.Wrapf(err, "can't create request to fetch pull request %s", id)
	}

	request.Header.Add("Authorization", fmt.Sprintf("token %s", g.Token))
	request.Header.Add("Accept", "application/vnd.github.v3+json")

	response, err := client.Do(request)
	if err != nil {
		return nil, errors.Wrapf(err, "can't fetch pull request %s", id)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "can't read github pull request %s response", id)
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("can't fetch pull request %s: %s", id, string(body))
	}

	var pullRequest PullRequestResponse

	err = json.Unmarshal(body, &pullRequest)
	if err != nil {
		return nil, errors.Wrapf(err, "can't parse github pull request %s response", id)
	}

	return &bugtracker.Issue{
		ID:      strconv.Itoa(pullRequest.ID),
		Subject: pullRequest.Subject,
		Link:    pullRequest.Link,
	}, nil
}

func githubPullRequestPath(github GitHub, id string) string {
	return fmt.Sprintf("%s/repos/%s/pulls/%s", github.API_URL, github.Repository, id)
}
