package changelog_test

import (
	"strings"
	"testing"
	"time"

	"github.com/kdisneur/changelog/pkg/changelog"
	"github.com/kdisneur/changelog/pkg/configuration"
	"github.com/kdisneur/changelog/pkg/formatter"
	"github.com/kdisneur/changelog/pkg/git"
	"github.com/kdisneur/changelog/pkg/github"
	"github.com/kdisneur/changelog/pkg/testing/bugtracker"
	"github.com/kdisneur/changelog/pkg/testing/repository"
)

func TestBuildChangelog(t *testing.T) {
	testCases := []struct {
		Name               string
		BuildConfiguration func() *configuration.ValidatedConfig
		IsValid            bool
		ErrorMessage       string
		ExpectedOutput     string
	}{
		{
			Name: "When configuration is valid",
			BuildConfiguration: func() *configuration.ValidatedConfig {
				tracker := bugtracker.NewBugTracker()
				repo := repository.New("git@github.com/kdisneur/changelog")

				repo.AddCommit(
					"7f76fa251d611ed48de62c460ec8f1b00804486b",
					git.Person{Fullname: "John Doe", Email: "john.doe@gmail.com"},
					time.Date(2018, time.November, 22, 5, 53, 12, 0, time.UTC),
					"initial Commit",
				)

				repo.AddCommit(
					"16dd9970c4f776157ccc6a7d8c78b2bdeeaab1c4",
					git.Person{Fullname: "John Doe", Email: "john.doe@gmail.com"},
					time.Date(2018, time.November, 22, 5, 56, 12, 0, time.UTC),
					"Add feature 1 (#1234)",
				)

				repo.AddCommit(
					"854da8029c41f552de16b81f7aba0e407a6bcb1c",
					git.Person{Fullname: "John Doe", Email: "john.doe@gmail.com"},
					time.Date(2018, time.November, 22, 5, 56, 12, 0, time.UTC),
					"Add feature 2 (#1337)",
				)

				tracker.AddIssue("1234", "Subject of feature 1")
				tracker.AddIssue("1337", "Subject of feature 2")

				return &configuration.ValidatedConfig{
					Repository:   repo,
					BugTracker:   tracker,
					From:         git.Reference("7f76fa251d611ed48de62c460ec8f1b00804486b"),
					To:           git.Reference("854da8029c41f552de16b81f7aba0e407a6bcb1c"),
					VersionName:  "v1.0.1",
					Date:         time.Date(2018, time.November, 22, 5, 59, 25, 0, time.UTC),
					CommitParser: github.NewSquashParser(),
					Formatter:    formatter.NewMarkdownFormatter(),
				}
			},
			IsValid:      true,
			ErrorMessage: "",
			ExpectedOutput: `## v1.0.1 - 2018-11-22

- Subject of feature 1 ([#1234])
- Subject of feature 2 ([#1337])

[#1234]: https://bugtracker.com/issue/1234
[#1337]: https://bugtracker.com/issue/1337
`,
		},
		{
			Name: "When there are no commits",
			BuildConfiguration: func() *configuration.ValidatedConfig {
				tracker := bugtracker.NewBugTracker()
				repo := repository.New("git@github.com/kdisneur/changelog")

				repo.AddCommit(
					"16dd9970c4f776157ccc6a7d8c78b2bdeeaab1c4",
					git.Person{Fullname: "John Doe", Email: "john.doe@gmail.com"},
					time.Date(2018, time.November, 22, 5, 56, 12, 0, time.UTC),
					"Add feature 1 (#1234)",
				)

				repo.AddCommit(
					"854da8029c41f552de16b81f7aba0e407a6bcb1c",
					git.Person{Fullname: "John Doe", Email: "john.doe@gmail.com"},
					time.Date(2018, time.November, 22, 5, 56, 12, 0, time.UTC),
					"Add feature 2 (#1337)",
				)

				tracker.AddIssue("1234", "Subject of feature 1")
				tracker.AddIssue("1337", "Subject of feature 2")

				return &configuration.ValidatedConfig{
					Repository:   repo,
					BugTracker:   tracker,
					From:         git.Reference("16dd9970c4f776157ccc6a7d8c78b2bdeeaab1c4"),
					To:           git.Reference("16dd9970c4f776157ccc6a7d8c78b2bdeeaab1c4"),
					VersionName:  "v1.0.1",
					Date:         time.Date(2018, time.November, 22, 5, 59, 25, 0, time.UTC),
					CommitParser: github.NewSquashParser(),
					Formatter:    formatter.NewMarkdownFormatter(),
				}
			},
			IsValid:        false,
			ErrorMessage:   "no commits found",
			ExpectedOutput: "",
		},
		{
			Name: "When log contains unparsable commits",
			BuildConfiguration: func() *configuration.ValidatedConfig {
				tracker := bugtracker.NewBugTracker()
				repo := repository.New("git@github.com/kdisneur/changelog")

				repo.AddCommit(
					"7f76fa251d611ed48de62c460ec8f1b00804486b",
					git.Person{Fullname: "John Doe", Email: "john.doe@gmail.com"},
					time.Date(2018, time.November, 22, 5, 53, 12, 0, time.UTC),
					"initial Commit",
				)

				repo.AddCommit(
					"16dd9970c4f776157ccc6a7d8c78b2bdeeaab1c4",
					git.Person{Fullname: "John Doe", Email: "john.doe@gmail.com"},
					time.Date(2018, time.November, 22, 5, 56, 12, 0, time.UTC),
					"Add feature 1 (#1234)",
				)

				repo.AddCommit(
					"854da8029c41f552de16b81f7aba0e407a6bcb1c",
					git.Person{Fullname: "John Doe", Email: "john.doe@gmail.com"},
					time.Date(2018, time.November, 22, 5, 56, 12, 0, time.UTC),
					"Add feature 2 #1337 with a wrong format",
				)

				tracker.AddIssue("1234", "Subject of feature 1")
				tracker.AddIssue("1337", "Subject of feature 2")

				return &configuration.ValidatedConfig{
					Repository:   repo,
					BugTracker:   tracker,
					From:         git.Reference("7f76fa251d611ed48de62c460ec8f1b00804486b"),
					To:           git.Reference("854da8029c41f552de16b81f7aba0e407a6bcb1c"),
					VersionName:  "v1.0.1",
					Date:         time.Date(2018, time.November, 22, 5, 59, 25, 0, time.UTC),
					CommitParser: github.NewSquashParser(),
					Formatter:    formatter.NewMarkdownFormatter(),
				}
			},
			IsValid:      true,
			ErrorMessage: "",
			ExpectedOutput: `## v1.0.1 - 2018-11-22

- Subject of feature 1 ([#1234])

[#1234]: https://bugtracker.com/issue/1234
`,
		},
		{
			Name: "When a commit contains a reference to a non existing issue",
			BuildConfiguration: func() *configuration.ValidatedConfig {
				tracker := bugtracker.NewBugTracker()
				repo := repository.New("git@github.com/kdisneur/changelog")

				repo.AddCommit(
					"7f76fa251d611ed48de62c460ec8f1b00804486b",
					git.Person{Fullname: "John Doe", Email: "john.doe@gmail.com"},
					time.Date(2018, time.November, 22, 5, 53, 12, 0, time.UTC),
					"initial Commit",
				)

				repo.AddCommit(
					"16dd9970c4f776157ccc6a7d8c78b2bdeeaab1c4",
					git.Person{Fullname: "John Doe", Email: "john.doe@gmail.com"},
					time.Date(2018, time.November, 22, 5, 56, 12, 0, time.UTC),
					"Add feature 1 (#1234)",
				)

				repo.AddCommit(
					"854da8029c41f552de16b81f7aba0e407a6bcb1c",
					git.Person{Fullname: "John Doe", Email: "john.doe@gmail.com"},
					time.Date(2018, time.November, 22, 5, 56, 12, 0, time.UTC),
					"Add feature 2 (#1337)",
				)

				tracker.AddIssue("42", "Subject of feature 1")
				tracker.AddIssue("1337", "Subject of feature 2")

				return &configuration.ValidatedConfig{
					Repository:   repo,
					BugTracker:   tracker,
					From:         git.Reference("7f76fa251d611ed48de62c460ec8f1b00804486b"),
					To:           git.Reference("854da8029c41f552de16b81f7aba0e407a6bcb1c"),
					VersionName:  "v1.0.1",
					Date:         time.Date(2018, time.November, 22, 5, 59, 25, 0, time.UTC),
					CommitParser: github.NewSquashParser(),
					Formatter:    formatter.NewMarkdownFormatter(),
				}
			},
			IsValid:        false,
			ErrorMessage:   "no issues with ID: 1234",
			ExpectedOutput: "",
		},
		{
			Name: "When no issues are found in matching commits",
			BuildConfiguration: func() *configuration.ValidatedConfig {
				tracker := bugtracker.NewBugTracker()
				repo := repository.New("git@github.com/kdisneur/changelog")

				repo.AddCommit(
					"7f76fa251d611ed48de62c460ec8f1b00804486b",
					git.Person{Fullname: "John Doe", Email: "john.doe@gmail.com"},
					time.Date(2018, time.November, 22, 5, 53, 12, 0, time.UTC),
					"initial Commit",
				)

				repo.AddCommit(
					"16dd9970c4f776157ccc6a7d8c78b2bdeeaab1c4",
					git.Person{Fullname: "John Doe", Email: "john.doe@gmail.com"},
					time.Date(2018, time.November, 22, 5, 56, 12, 0, time.UTC),
					"Add feature 1",
				)

				repo.AddCommit(
					"854da8029c41f552de16b81f7aba0e407a6bcb1c",
					git.Person{Fullname: "John Doe", Email: "john.doe@gmail.com"},
					time.Date(2018, time.November, 22, 5, 56, 12, 0, time.UTC),
					"Add feature 2",
				)

				tracker.AddIssue("42", "Subject of feature 1")
				tracker.AddIssue("1337", "Subject of feature 2")

				return &configuration.ValidatedConfig{
					Repository:   repo,
					BugTracker:   tracker,
					From:         git.Reference("7f76fa251d611ed48de62c460ec8f1b00804486b"),
					To:           git.Reference("854da8029c41f552de16b81f7aba0e407a6bcb1c"),
					VersionName:  "v1.0.1",
					Date:         time.Date(2018, time.November, 22, 5, 59, 25, 0, time.UTC),
					CommitParser: github.NewSquashParser(),
					Formatter:    formatter.NewMarkdownFormatter(),
				}
			},
			IsValid:        false,
			ErrorMessage:   "no commits kept",
			ExpectedOutput: "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			config := testCase.BuildConfiguration()

			output, err := changelog.BuildChangelog(config)

			if err != nil && testCase.IsValid {
				t.Fatalf("Expected no errors but got: %s", err.Error())
			}

			if err == nil && !testCase.IsValid {
				t.Fatalf("Expected errors but got none.\nOutput\n %s", output)
			}

			if !testCase.IsValid && !strings.Contains(err.Error(), testCase.ErrorMessage) {
				t.Fatalf("Wrong error. Expected: %s\nReceived: %s", testCase.ErrorMessage, err.Error())
			}

			if output != testCase.ExpectedOutput {
				t.Fatalf("Wrong output. Expected:\n%s\nReceived:\n%s", testCase.ExpectedOutput, output)
			}
		})
	}
}
