package bugtracker

import (
	"fmt"
	"github.com/kdisneur/changelog/pkg/bugtracker"
)

type BugTracker struct {
	Issues map[string]*bugtracker.Issue
}

func NewBugTracker() *BugTracker {
	return &BugTracker{
		Issues: make(map[string]*bugtracker.Issue),
	}
}

func (b BugTracker) Equal(other bugtracker.BugTracker) bool {
	_, hasGoodType := other.(BugTracker)

	return hasGoodType
}

func (b BugTracker) FindIssue(id string) (*bugtracker.Issue, error) {
	issue := b.Issues[id]

	if issue == nil {
		return nil, fmt.Errorf("no issues with ID: %s", id)
	}

	return issue, nil
}

func (b *BugTracker) AddIssue(id string, subject string) {
	b.Issues[id] = &bugtracker.Issue{
		ID:      id,
		Subject: subject,
		Link:    fmt.Sprintf("https://bugtracker.com/issue/%s", id),
	}
}
