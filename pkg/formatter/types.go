package formatter

import (
	"github.com/kdisneur/changelog/pkg/bugtracker"
	"time"
)

type Formatter interface {
	Equal(other Formatter) bool
	FormatIssues(versionName string, date time.Time, issues []*bugtracker.Issue) string
}
