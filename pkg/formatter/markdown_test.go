package formatter_test

import (
	"testing"
	"time"

	"github.com/kdisneur/changelog/pkg/bugtracker"
	"github.com/kdisneur/changelog/pkg/formatter"
)

func TestMardownFormatter(t *testing.T) {
	testCases := []struct {
		Name     string
		Version  string
		Date     time.Time
		Issues   []*bugtracker.Issue
		Expected string
	}{
		{
			"When it contains several issues",
			"v1.0.0",
			time.Date(2018, time.November, 19, 5, 12, 42, 0, time.UTC),
			[]*bugtracker.Issue{
				&bugtracker.Issue{ID: "42", Subject: "A nice feature", Link: "https://github.com/kdisneur/changelog/pull/42"},
				&bugtracker.Issue{ID: "1337", Subject: "Another nice feature", Link: "https://github.com/kdisneur/changelog/pull/1337"},
			},
			`## v1.0.0 - 2018-11-19

- A nice feature ([#42])
- Another nice feature ([#1337])

[#42]: https://github.com/kdisneur/changelog/pull/42
[#1337]: https://github.com/kdisneur/changelog/pull/1337
`,
		},
		{
			"When it contains one issue",
			"v1.0.0",
			time.Date(2018, time.November, 19, 5, 12, 42, 0, time.UTC),
			[]*bugtracker.Issue{
				&bugtracker.Issue{ID: "42", Subject: "A nice feature", Link: "https://github.com/kdisneur/changelog/pull/42"},
			},
			`## v1.0.0 - 2018-11-19

- A nice feature ([#42])

[#42]: https://github.com/kdisneur/changelog/pull/42
`,
		},
		{
			"When it contains no issues",
			"v1.0.0",
			time.Date(2018, time.November, 19, 5, 12, 42, 0, time.UTC),
			[]*bugtracker.Issue{},
			`## v1.0.0 - 2018-11-19

(No changes)
`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			markdown := formatter.NewMarkdownFormatter()
			actual := markdown.FormatIssues(
				testCase.Version,
				testCase.Date,
				testCase.Issues,
			)

			if actual != testCase.Expected {
				t.Errorf("Wrong format.\nExpected:\n%s\n\nReceived:\n%s", testCase.Expected, actual)
			}
		})
	}
}
