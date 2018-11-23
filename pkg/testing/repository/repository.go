package repository

import (
	"github.com/kdisneur/changelog/pkg/git"
	"github.com/kdisneur/changelog/pkg/git/utils"
	"time"
)

type Repository struct {
	remoteURL string
	Commits   []*git.Commit
}

func New(remoteURL string) *Repository {
	return &Repository{remoteURL: remoteURL}
}

func (r Repository) Equal(other git.Git) bool {
	repository, hasGoodType := other.(Repository)
	if !hasGoodType {
		return false
	}

	return r.remoteURL == repository.remoteURL
}

func (r Repository) Log(from git.Reference, to git.Reference) ([]*git.Commit, error) {
	var commits []*git.Commit
	shouldKeep := false

	for _, commit := range r.Commits {
		if shouldKeep {
			commits = append(commits, commit)
		}

		if commit.ID == string(to) {
			break
		}

		if commit.ID == string(from) {
			shouldKeep = true
		}

	}

	return commits, nil
}

func (r Repository) FindRemote() (*git.Remote, error) {
	return utils.FindRemoteFromURLs([]string{r.remoteURL})
}

func (r *Repository) AddCommit(id string, author git.Person, authoredAt time.Time, message string) {
	r.Commits = append(r.Commits, buildCommit(id, author, authoredAt, message, false))
}

func (r *Repository) AddMergeCommit(id string, author git.Person, authoredAt time.Time, message string) {
	r.Commits = append(r.Commits, buildCommit(id, author, authoredAt, message, true))
}

func buildCommit(id string, author git.Person, authoredAt time.Time, message string, merge bool) *git.Commit {
	return &git.Commit{
		ID:          id,
		Author:      author,
		AuthoredAt:  authoredAt,
		Committer:   author,
		CommittedAt: authoredAt,
		IsMerge:     merge,
		Message:     message,
	}
}
