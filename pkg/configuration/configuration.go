package configuration

import (
	"fmt"

	"github.com/kdisneur/changelog/pkg/formatter"
	"github.com/kdisneur/changelog/pkg/git"
	"github.com/kdisneur/changelog/pkg/git/system"
	"github.com/kdisneur/changelog/pkg/github"
	"github.com/kdisneur/changelog/pkg/parser"
)

func Validate(file File, command Command) (*ValidatedConfig, error) {
	repository, err := system.NewRepository(command.RepositoryLocalPath)
	if err != nil {
		return nil, err
	}

	repositoryName, err := getRepositoryName(repository, command)
	if err != nil {
		return nil, err
	}

	fromReference := git.NewReference(command.From)
	toReference := getToReference(file, command, repositoryName)

	commitParser, err := getCommitParser(file, command, repositoryName)
	if err != nil {
		return nil, err
	}

	formatter := formatter.NewMarkdownFormatter()

	tracker := github.NewBugTracker(file.Github.Token, repositoryName)

	return &ValidatedConfig{
		From:         fromReference,
		To:           toReference,
		VersionName:  command.VersionName,
		Date:         command.Date,
		CommitParser: commitParser,
		Formatter:    formatter,
		Repository:   repository,
		BugTracker:   tracker,
	}, nil
}

func getCommitParser(file File, command Command, repositoryName string) (parser.Parser, error) {
	strategy := "squash"

	if command.MergeStrategy != "" {
		strategy = command.MergeStrategy
	} else {
		repository, ok := file.FindRepository(repositoryName)
		if ok && repository.MergeStrategy != "" {
			strategy = repository.MergeStrategy
		} else if file.General.MergeStrategy != "" {
			strategy = file.General.MergeStrategy
		}
	}

	switch strategy {
	case "squash":
		return github.NewSquashParser(), nil
	case "merge":
		return github.NewMergeParser(), nil
	default:
		return nil, fmt.Errorf("Asked for '%s' strategy but support only 'squash' and 'merge'", strategy)
	}
}

func getToReference(file File, command Command, repositoryName string) git.Reference {
	if command.To != "" {
		return git.NewReference(command.To)
	}

	repository, ok := file.FindRepository(repositoryName)
	if ok && repository.BaseBranch != "" {
		return git.NewReference(repository.BaseBranch)
	}

	if file.General.BaseBranch != "" {
		return git.NewReference(file.General.BaseBranch)
	}

	return git.NewReference("master")
}

func getRepositoryName(repository git.Git, command Command) (string, error) {
	repositoryName := command.RepositoryName
	if repositoryName == "" {
		remote, err := repository.FindRemote()
		if err != nil {
			return "", err
		}
		repositoryName = remote.RepositoryName
	}

	return repositoryName, nil
}
