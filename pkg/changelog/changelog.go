package changelog

import (
	"errors"

	"github.com/kdisneur/changelog/pkg/bugtracker"
	"github.com/kdisneur/changelog/pkg/configuration"
)

func BuildChangelog(conf *configuration.ValidatedConfig) (string, error) {
	commits, err := conf.Repository.Log(conf.From, conf.To)
	if err != nil {
		return "", err
	}

	if len(commits) == 0 {
		return "", errors.New("no commits found")
	}

	var issues []*bugtracker.Issue
	for _, commit := range commits {
		if conf.CommitParser.KeepCommit(commit.Message) {
			id, err := conf.CommitParser.FindID(commit.Message)
			if err != nil {
				return "", err
			}

			issue, err := conf.BugTracker.FindIssue(id)
			if err != nil {
				return "", err
			}

			issues = append(issues, issue)
		}
	}

	if len(issues) == 0 {
		return "", errors.New("no commits kept")
	}

	return conf.Formatter.FormatIssues(conf.VersionName, conf.Date, issues), nil
}
