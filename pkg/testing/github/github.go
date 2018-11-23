package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const body string = `
A long body describing
the feature just added
`

func NewMock(token string, repository string, id int, number int, title string) GitHubMock {
	return GitHubMock{
		Token:      token,
		Repository: repository,
		PullRequest: PullRequest{
			ID:       id,
			Number:   number,
			URL:      fmt.Sprintf("https://api.github.com/repos/%s/pulls/%d", repository, number),
			HTML_URL: fmt.Sprintf("https://github.com/%s/pulls/%d", repository, number),
			IssueURL: fmt.Sprintf("https://github.com/%s/issues/%d", repository, number),
			Title:    title,
			Body:     body,
		},
	}
}

func (m *GitHubMock) Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	pullRequestURL, err := url.ParseRequestURI(m.PullRequest.URL)
	if err != nil || r.URL.Path != pullRequestURL.Path {
		response, _ := json.Marshal(HTTPError{"Not Found", "https://developer.github.com/v3/pulls/#get-a-single-pull-request"})
		http.Error(w, string(response), 404)

		return
	}

	token := strings.TrimPrefix(r.Header.Get("Authorization"), "token ")
	if m.Token != token {
		response, _ := json.Marshal(HTTPError{"Bad credentials", "https://developer.github.com"})
		http.Error(w, string(response), 401)

		return
	}

	response, _ := json.Marshal(m.PullRequest)
	w.Write(response)
}
