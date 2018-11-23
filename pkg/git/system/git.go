package system

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/kdisneur/changelog/pkg/git"
	sysutils "github.com/kdisneur/changelog/pkg/git/system/utils"
	"github.com/kdisneur/changelog/pkg/git/utils"
	"github.com/kdisneur/changelog/pkg/time"
	"github.com/pkg/errors"
)

type Repository struct {
	RepositoryPath git.Path
}

func NewRepository(repositoryPath string) (git.Git, error) {
	if !sysutils.IsGitRepository(repositoryPath) {
		return nil, errors.New(fmt.Sprintf("Path '%s' is not a git repository", repositoryPath))
	}

	return &Repository{RepositoryPath: git.Path(repositoryPath)}, nil
}

func (r Repository) Equal(other git.Git) bool {
	otherRepository, hasGoodType := other.(*Repository)
	if !hasGoodType {
		return false
	}

	return r.RepositoryPath == otherRepository.RepositoryPath
}

func (r Repository) FindRemote() (*git.Remote, error) {
	commandResult, err := sysutils.ExecCommand(r.RepositoryPath.String(), "remote", "--verbose")

	if err != nil {
		return nil, errors.Wrapf(err, "Can't find git remotes")
	}

	rawRemotes := strings.Split(commandResult, "\n")

	var remoteURLs []string
	for _, rawRemote := range rawRemotes {
		remoteData := strings.Fields(rawRemote)
		if len(remoteData) > 1 {
			remoteURLs = append(remoteURLs, remoteData[1])
		}
	}

	return utils.FindRemoteFromURLs(remoteURLs)
}

func (r Repository) Log(from git.Reference, to git.Reference) ([]*git.Commit, error) {
	span := fmt.Sprintf("%s..%s", string(from), string(to))

	rawCommits, err := sysutils.ExecCommand(r.RepositoryPath.String(), "log", "--pretty=oneline", "--format=%H;%an;%aE;%at;%cn;%ce;%ct;%p;%s", span)

	if err != nil {
		return nil, errors.Wrapf(err, "Can't generate git logs for '%s' in %s", span, r.RepositoryPath)
	}

	return parseRawCommits(rawCommits)
}

func parseRawCommits(rawCommit string) ([]*git.Commit, error) {
	scanner := bufio.NewScanner(strings.NewReader(rawCommit))

	var commits []*git.Commit
	for scanner.Scan() {
		commit, err := parseRawCommit(scanner.Text())
		if err != nil {
			return nil, err
		}

		commits = append(commits, commit)
	}

	return commits, nil
}

func parseRawCommit(rawCommit string) (*git.Commit, error) {
	commitData := strings.Split(rawCommit, ";")
	if len(commitData) != 9 {
		return nil, errors.New(fmt.Sprintf("Can't parse '%s'", rawCommit))
	}

	id := commitData[0]
	authorName := commitData[1]
	authorEmail := commitData[2]
	authorTimestamp := commitData[3]
	committerName := commitData[4]
	committerEmail := commitData[5]
	commitTimestamp := commitData[6]
	parentHashes := commitData[7]
	message := commitData[8]

	author := git.NewPerson(authorName, authorEmail)
	committer := git.NewPerson(committerName, committerEmail)
	authoredAt, err := time.FromStringTimestamp(authorTimestamp)
	if err != nil {
		return nil, errors.Wrap(err, "Can't parse author timestamp")
	}

	committedAt, err := time.FromStringTimestamp(commitTimestamp)
	if err != nil {
		return nil, errors.Wrap(err, "Can't parse committer timestamp")
	}

	isMerge := len(strings.Split(parentHashes, " ")) > 1

	return &git.Commit{
		ID:          id,
		Author:      author,
		AuthoredAt:  authoredAt,
		Committer:   committer,
		CommittedAt: committedAt,
		IsMerge:     isMerge,
		Message:     message,
	}, nil
}
