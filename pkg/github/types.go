package github

type GitHub struct {
	Token      string
	API_URL    string
	Repository string
}

type PullRequestResponse struct {
	ID      int    `json:"number"`
	Subject string `json:"title"`
	Link    string `json:"html_url"`
}
