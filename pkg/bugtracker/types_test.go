package bugtracker_test

import (
	"testing"

	"github.com/kdisneur/changelog/pkg/bugtracker"
)

func TestIssueEqual(t *testing.T) {
	testCases := []struct {
		Name     string
		Issue1   bugtracker.Issue
		Issue2   bugtracker.Issue
		Expected bool
	}{
		{
			"When issues are equal",
			bugtracker.Issue{ID: "42", Subject: "A good subject name", Link: "https://site.com/issue/42"},
			bugtracker.Issue{ID: "42", Subject: "A good subject name", Link: "https://site.com/issue/42"},
			true,
		},
		{
			"When issues have different IDs",
			bugtracker.Issue{ID: "42", Subject: "A good subject name", Link: "https://site.com/issue/42"},
			bugtracker.Issue{ID: "1337", Subject: "A good subject name", Link: "https://site.com/issue/42"},
			false,
		},
		{
			"When issues have different subjects",
			bugtracker.Issue{ID: "42", Subject: "A good subject name", Link: "https://site.com/issue/42"},
			bugtracker.Issue{ID: "42", Subject: "Another subject name", Link: "https://site.com/issue/42"},
			false,
		}, {
			"When issues have different link URLs",
			bugtracker.Issue{ID: "42", Subject: "A good subject name", Link: "https://site.com/issue/42"},
			bugtracker.Issue{ID: "42", Subject: "A good subject name", Link: "https://anothersite.com/issue/42"},
			false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			word := "equality"
			if !testCase.Expected {
				word = "inequality"
			}

			if testCase.Issue1.Equal(&testCase.Issue2) != testCase.Expected {
				t.Errorf("Expected %s.\nIssue1:\n%+v\nIssue2:\n%+v", word, testCase.Issue1, testCase.Issue2)
			}
		})
	}
}
