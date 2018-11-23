package github

type PullRequest struct {
	ID       int
	Number   int
	URL      string
	HTML_URL string
	IssueURL string
	Title    string
	Body     string
}

type HTTPError struct {
	Message          string `json:"message"`
	DocumentationURL string `json:"documentation_url"`
}

type GitHubMock struct {
	Token       string
	Repository  string
	PullRequest PullRequest
}
