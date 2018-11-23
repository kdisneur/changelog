package configuration

import (
	"time"

	"github.com/kdisneur/changelog/pkg/bugtracker"
	"github.com/kdisneur/changelog/pkg/formatter"
	"github.com/kdisneur/changelog/pkg/git"
	"github.com/kdisneur/changelog/pkg/parser"
)

type File struct {
	General    General
	Github     GitHub
	Repository []GitRepository
}

type General struct {
	MergeStrategy string
	BaseBranch    string
}

type GitHub struct {
	Token string
}

type GitRepository struct {
	Name          string
	BaseBranch    string
	MergeStrategy string
}

type Command struct {
	RepositoryName      string
	From                string
	To                  string
	VersionName         string
	Date                time.Time
	RepositoryLocalPath string
	MergeStrategy       string
}

type ValidatedConfig struct {
	From         git.Reference
	To           git.Reference
	VersionName  string
	Date         time.Time
	CommitParser parser.Parser
	Formatter    formatter.Formatter
	Repository   git.Git
	BugTracker   bugtracker.BugTracker
}

func (c *ValidatedConfig) Equal(other *ValidatedConfig) bool {
	return c.From == other.From &&
		c.To == other.To &&
		c.VersionName == other.VersionName &&
		c.Date.Equal(other.Date) &&
		c.CommitParser.Equal(other.CommitParser) &&
		c.Formatter.Equal(other.Formatter) &&
		c.Repository.Equal(other.Repository) &&
		c.BugTracker.Equal(other.BugTracker)
}
