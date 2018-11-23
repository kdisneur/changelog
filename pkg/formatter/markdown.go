package formatter

import (
	"bytes"
	"fmt"
	"time"

	"github.com/kdisneur/changelog/pkg/bugtracker"
)

type markdownFormatter struct{}

func NewMarkdownFormatter() Formatter {
	return &markdownFormatter{}
}

func (m markdownFormatter) Equal(other Formatter) bool {
	_, hasGoodType := other.(*markdownFormatter)

	return hasGoodType
}

func (m markdownFormatter) FormatIssues(versionName string, releaseDate time.Time, issues []*bugtracker.Issue) string {
	if len(issues) == 0 {
		return formatNoIssues(versionName, releaseDate)
	} else {
		return formatIssues(versionName, releaseDate, issues)
	}
}

func formatIssues(versionName string, releaseDate time.Time, issues []*bugtracker.Issue) string {
	var list bytes.Buffer
	var links bytes.Buffer

	for _, issue := range issues {
		list.WriteString(fmt.Sprintf("- %s ([#%s])\n", issue.Subject, issue.ID))
		links.WriteString(fmt.Sprintf("[#%s]: %s\n", issue.ID, issue.Link))
	}

	return fmt.Sprintf("## %s - %s\n\n%s\n%s", versionName, formatReleaseDate(releaseDate), list.String(), links.String())
}
func formatNoIssues(versionName string, releaseDate time.Time) string {
	return fmt.Sprintf("## %s - %s\n\n(No changes)\n", versionName, formatReleaseDate(releaseDate))
}

func formatReleaseDate(date time.Time) string {
	return date.Format("2006-01-02")
}
