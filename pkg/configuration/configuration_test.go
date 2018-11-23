package configuration_test

import (
	"strings"
	"testing"
	"time"

	"github.com/kdisneur/changelog/pkg/configuration"
	"github.com/kdisneur/changelog/pkg/formatter"
	"github.com/kdisneur/changelog/pkg/git"
	"github.com/kdisneur/changelog/pkg/git/system"
	"github.com/kdisneur/changelog/pkg/github"
	"github.com/kdisneur/changelog/pkg/testing/targz"
)

func TestValidate(t *testing.T) {
	const ValidGitHubToken string = "aaaa-bbbb-cccc-dddd"
	const ValidRepositoryName string = "kdisneur/changelog"

	testCases := []struct {
		Name                       string
		File                       configuration.File
		CommandRepositoryName      string
		CommandFrom                string
		CommandTo                  string
		CommandVersionName         string
		CommandDate                time.Time
		CommandRepositoryLocalPath string
		CommandMergeStrategy       string
		Fixture                    string
		IsValid                    bool
		ErrorMessage               string
		ExpectedBuilder            func(path string) *configuration.ValidatedConfig
	}{
		{
			Name: "When configuration is full",
			File: configuration.File{
				Github: configuration.GitHub{Token: ValidGitHubToken},
			},
			CommandRepositoryName: ValidRepositoryName,
			CommandFrom:           "v1.0.0",
			CommandTo:             "master",
			CommandVersionName:    "v1.0.1",
			CommandDate:           time.Date(2018, time.November, 21, 5, 45, 12, 0, time.UTC),
			CommandMergeStrategy:  "squash",
			Fixture:               "squash",
			IsValid:               true,
			ErrorMessage:          "",
			ExpectedBuilder: func(path string) *configuration.ValidatedConfig {
				repository, _ := system.NewRepository(path)

				return &configuration.ValidatedConfig{
					From:         git.Reference("v1.0.0"),
					To:           git.Reference("master"),
					VersionName:  "v1.0.1",
					Date:         time.Date(2018, time.November, 21, 5, 45, 12, 0, time.UTC),
					CommitParser: github.NewSquashParser(),
					Formatter:    formatter.NewMarkdownFormatter(),
					Repository:   repository,
					BugTracker:   github.NewBugTracker(ValidGitHubToken, ValidRepositoryName),
				}
			},
		},
		{
			Name: "When configuration has no command branch definition but a file one",
			File: configuration.File{
				Github: configuration.GitHub{Token: ValidGitHubToken},
				Repository: []configuration.GitRepository{
					{
						Name:       "another/repository",
						BaseBranch: "master",
					},
					{
						Name:       ValidRepositoryName,
						BaseBranch: "develop",
					},
				},
			},
			CommandRepositoryName: ValidRepositoryName,
			CommandFrom:           "v1.0.0",
			CommandTo:             "",
			CommandVersionName:    "v1.0.1",
			CommandDate:           time.Date(2018, time.November, 21, 5, 45, 12, 0, time.UTC),
			CommandMergeStrategy:  "squash",
			Fixture:               "squash",
			IsValid:               true,
			ErrorMessage:          "",
			ExpectedBuilder: func(path string) *configuration.ValidatedConfig {
				repository, _ := system.NewRepository(path)

				return &configuration.ValidatedConfig{
					From:         git.Reference("v1.0.0"),
					To:           git.Reference("develop"),
					VersionName:  "v1.0.1",
					Date:         time.Date(2018, time.November, 21, 5, 45, 12, 0, time.UTC),
					CommitParser: github.NewSquashParser(),
					Formatter:    formatter.NewMarkdownFormatter(),
					Repository:   repository,
					BugTracker:   github.NewBugTracker(ValidGitHubToken, ValidRepositoryName),
				}
			},
		},
		{
			Name: "When configuration has command and file branch definition",
			File: configuration.File{
				Github: configuration.GitHub{Token: ValidGitHubToken},
				Repository: []configuration.GitRepository{
					{
						Name:       "another/repository",
						BaseBranch: "master",
					},
					{
						Name:       ValidRepositoryName,
						BaseBranch: "fromfile",
					},
				},
			},
			CommandRepositoryName: ValidRepositoryName,
			CommandFrom:           "v1.0.0",
			CommandTo:             "fromcommand",
			CommandVersionName:    "v1.0.1",
			CommandDate:           time.Date(2018, time.November, 21, 5, 45, 12, 0, time.UTC),
			CommandMergeStrategy:  "squash",
			Fixture:               "squash",
			IsValid:               true,
			ErrorMessage:          "",
			ExpectedBuilder: func(path string) *configuration.ValidatedConfig {
				repository, _ := system.NewRepository(path)

				return &configuration.ValidatedConfig{
					From:         git.Reference("v1.0.0"),
					To:           git.Reference("fromcommand"),
					VersionName:  "v1.0.1",
					Date:         time.Date(2018, time.November, 21, 5, 45, 12, 0, time.UTC),
					CommitParser: github.NewSquashParser(),
					Formatter:    formatter.NewMarkdownFormatter(),
					Repository:   repository,
					BugTracker:   github.NewBugTracker(ValidGitHubToken, ValidRepositoryName),
				}
			},
		},
		{
			Name: "When configuration has no command nor file branch definition",
			File: configuration.File{
				Github: configuration.GitHub{Token: ValidGitHubToken},
				Repository: []configuration.GitRepository{
					{
						Name:       "another/repository",
						BaseBranch: "develop",
					},
					{
						Name: ValidRepositoryName,
					},
				},
			},
			CommandRepositoryName: ValidRepositoryName,
			CommandFrom:           "v1.0.0",
			CommandTo:             "",
			CommandVersionName:    "v1.0.1",
			CommandDate:           time.Date(2018, time.November, 21, 5, 45, 12, 0, time.UTC),
			CommandMergeStrategy:  "squash",
			Fixture:               "squash",
			IsValid:               true,
			ErrorMessage:          "",
			ExpectedBuilder: func(path string) *configuration.ValidatedConfig {
				repository, _ := system.NewRepository(path)

				return &configuration.ValidatedConfig{
					From:         git.Reference("v1.0.0"),
					To:           git.Reference("master"),
					VersionName:  "v1.0.1",
					Date:         time.Date(2018, time.November, 21, 5, 45, 12, 0, time.UTC),
					CommitParser: github.NewSquashParser(),
					Formatter:    formatter.NewMarkdownFormatter(),
					Repository:   repository,
					BugTracker:   github.NewBugTracker(ValidGitHubToken, ValidRepositoryName),
				}
			},
		},
		{
			Name: "When configuration has no command merge strategy definition but a file one",
			File: configuration.File{
				Github: configuration.GitHub{Token: ValidGitHubToken},
				Repository: []configuration.GitRepository{
					{
						Name:       "another/repository",
						BaseBranch: "master",
					},
					{
						Name:          ValidRepositoryName,
						MergeStrategy: "merge",
					},
				},
			},
			CommandRepositoryName: ValidRepositoryName,
			CommandFrom:           "v1.0.0",
			CommandTo:             "master",
			CommandVersionName:    "v1.0.1",
			CommandDate:           time.Date(2018, time.November, 21, 5, 45, 12, 0, time.UTC),
			CommandMergeStrategy:  "",
			Fixture:               "squash",
			IsValid:               true,
			ErrorMessage:          "",
			ExpectedBuilder: func(path string) *configuration.ValidatedConfig {
				repository, _ := system.NewRepository(path)

				return &configuration.ValidatedConfig{
					From:         git.Reference("v1.0.0"),
					To:           git.Reference("master"),
					VersionName:  "v1.0.1",
					Date:         time.Date(2018, time.November, 21, 5, 45, 12, 0, time.UTC),
					CommitParser: github.NewMergeParser(),
					Formatter:    formatter.NewMarkdownFormatter(),
					Repository:   repository,
					BugTracker:   github.NewBugTracker(ValidGitHubToken, ValidRepositoryName),
				}
			},
		},
		{
			Name: "When configuration has command and file merge strategy definition",
			File: configuration.File{
				Github: configuration.GitHub{Token: ValidGitHubToken},
				Repository: []configuration.GitRepository{
					{
						Name:       "another/repository",
						BaseBranch: "master",
					},
					{
						Name:          ValidRepositoryName,
						MergeStrategy: "merge",
					},
				},
			},
			CommandRepositoryName: ValidRepositoryName,
			CommandFrom:           "v1.0.0",
			CommandTo:             "master",
			CommandVersionName:    "v1.0.1",
			CommandDate:           time.Date(2018, time.November, 21, 5, 45, 12, 0, time.UTC),
			CommandMergeStrategy:  "squash",
			Fixture:               "squash",
			IsValid:               true,
			ErrorMessage:          "",
			ExpectedBuilder: func(path string) *configuration.ValidatedConfig {
				repository, _ := system.NewRepository(path)

				return &configuration.ValidatedConfig{
					From:         git.Reference("v1.0.0"),
					To:           git.Reference("master"),
					VersionName:  "v1.0.1",
					Date:         time.Date(2018, time.November, 21, 5, 45, 12, 0, time.UTC),
					CommitParser: github.NewSquashParser(),
					Formatter:    formatter.NewMarkdownFormatter(),
					Repository:   repository,
					BugTracker:   github.NewBugTracker(ValidGitHubToken, ValidRepositoryName),
				}
			},
		},
		{
			Name: "When configuration has no command nor file merge strategy definition",
			File: configuration.File{
				Github: configuration.GitHub{Token: ValidGitHubToken},
				Repository: []configuration.GitRepository{
					{
						Name:       "another/repository",
						BaseBranch: "develop",
					},
					{
						Name: ValidRepositoryName,
					},
				},
			},
			CommandRepositoryName: ValidRepositoryName,
			CommandFrom:           "v1.0.0",
			CommandTo:             "",
			CommandVersionName:    "v1.0.1",
			CommandDate:           time.Date(2018, time.November, 21, 5, 45, 12, 0, time.UTC),
			CommandMergeStrategy:  "",
			Fixture:               "squash",
			IsValid:               true,
			ErrorMessage:          "",
			ExpectedBuilder: func(path string) *configuration.ValidatedConfig {
				repository, _ := system.NewRepository(path)

				return &configuration.ValidatedConfig{
					From:         git.Reference("v1.0.0"),
					To:           git.Reference("master"),
					VersionName:  "v1.0.1",
					Date:         time.Date(2018, time.November, 21, 5, 45, 12, 0, time.UTC),
					CommitParser: github.NewSquashParser(),
					Formatter:    formatter.NewMarkdownFormatter(),
					Repository:   repository,
					BugTracker:   github.NewBugTracker(ValidGitHubToken, ValidRepositoryName),
				}
			},
		},
		{
			Name: "When configuration has no repository name but one remote is defined in Git",
			File: configuration.File{
				Github: configuration.GitHub{Token: ValidGitHubToken},
			},
			CommandRepositoryName: "",
			CommandFrom:           "v1.0.0",
			CommandTo:             "master",
			CommandVersionName:    "v1.0.1",
			CommandDate:           time.Date(2018, time.November, 21, 5, 45, 12, 0, time.UTC),
			CommandMergeStrategy:  "squash",
			Fixture:               "onegitremote",
			IsValid:               true,
			ErrorMessage:          "",
			ExpectedBuilder: func(path string) *configuration.ValidatedConfig {
				repository, _ := system.NewRepository(path)

				return &configuration.ValidatedConfig{
					From:         git.Reference("v1.0.0"),
					To:           git.Reference("master"),
					VersionName:  "v1.0.1",
					Date:         time.Date(2018, time.November, 21, 5, 45, 12, 0, time.UTC),
					CommitParser: github.NewSquashParser(),
					Formatter:    formatter.NewMarkdownFormatter(),
					Repository:   repository,
					BugTracker:   github.NewBugTracker(ValidGitHubToken, ValidRepositoryName),
				}
			},
		},
		{
			Name: "When configuration has no repository name but multiple remotes are defined in Git",
			File: configuration.File{
				Github: configuration.GitHub{Token: ValidGitHubToken},
			},
			CommandRepositoryName: "",
			CommandFrom:           "v1.0.0",
			CommandTo:             "master",
			CommandVersionName:    "v1.0.1",
			CommandDate:           time.Date(2018, time.November, 21, 5, 45, 12, 0, time.UTC),
			CommandMergeStrategy:  "squash",
			Fixture:               "multipledifferentremotes",
			IsValid:               false,
			ErrorMessage:          "found multiple remotes",
			ExpectedBuilder:       func(path string) *configuration.ValidatedConfig { return nil },
		},
		{
			Name: "When configuration contains an unsupported merging strategy",
			File: configuration.File{
				Github: configuration.GitHub{Token: ValidGitHubToken},
			},
			CommandRepositoryName: ValidRepositoryName,
			CommandFrom:           "v1.0.0",
			CommandTo:             "master",
			CommandVersionName:    "v1.0.1",
			CommandDate:           time.Date(2018, time.November, 21, 5, 45, 12, 0, time.UTC),
			CommandMergeStrategy:  "wrong-strategy",
			Fixture:               "squash",
			IsValid:               false,
			ErrorMessage:          "Asked for 'wrong-strategy' strategy but support only",
			ExpectedBuilder:       func(path string) *configuration.ValidatedConfig { return nil },
		},
		{
			Name: "When configuration contains a path to a wrong repository",
			File: configuration.File{
				Github: configuration.GitHub{Token: ValidGitHubToken},
			},
			CommandRepositoryName:      ValidRepositoryName,
			CommandRepositoryLocalPath: "/tmp/not/git/repo",
			CommandFrom:                "v1.0.0",
			CommandTo:                  "master",
			CommandVersionName:         "v1.0.1",
			CommandDate:                time.Date(2018, time.November, 21, 5, 45, 12, 0, time.UTC),
			CommandMergeStrategy:       "squash",
			Fixture:                    "squash",
			IsValid:                    false,
			ErrorMessage:               "Path '/tmp/not/git/repo' is not a git repository",
			ExpectedBuilder:            func(path string) *configuration.ValidatedConfig { return nil },
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			if testCase.CommandRepositoryLocalPath == "" && testCase.Fixture != "" {
				repositoryPath, cleanup, err := targz.Untar(testCase.Fixture)
				if err != nil {
					t.Fatalf(err.Error())
				}

				defer cleanup()

				testCase.CommandRepositoryLocalPath = repositoryPath
			}

			command := configuration.Command{
				RepositoryName:      testCase.CommandRepositoryName,
				From:                testCase.CommandFrom,
				To:                  testCase.CommandTo,
				VersionName:         testCase.CommandVersionName,
				Date:                testCase.CommandDate,
				RepositoryLocalPath: testCase.CommandRepositoryLocalPath,
				MergeStrategy:       testCase.CommandMergeStrategy,
			}

			config, err := configuration.Validate(testCase.File, command)

			if err != nil && testCase.IsValid {
				t.Fatalf("Expected no errors, got one. Received: %s", err.Error())
			}

			if err == nil && !testCase.IsValid {
				t.Fatalf("Expected an error, got none. Received: %+v", config)
			}

			if !testCase.IsValid && !strings.Contains(err.Error(), testCase.ErrorMessage) {
				t.Fatalf("Wrong error.\nExpected: %s\nReceived: %s", testCase.ErrorMessage, err.Error())
			}

			if testCase.IsValid {
				expected := testCase.ExpectedBuilder(testCase.CommandRepositoryLocalPath)
				if !expected.Equal(config) {
					t.Fatalf("Wrong configuration.\nExpected:\n%+v\nReceived:\n%+v", expected, config)
				}
			}
		})
	}
}
