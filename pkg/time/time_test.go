package time_test

import (
	"strings"
	"testing"
	"time"

	timeutils "github.com/kdisneur/changelog/pkg/time"
)

func TestStringTimestamp(t *testing.T) {
	secondsCETOfUTC := int((1 * time.Hour).Seconds())
	centralEuropeTime := time.FixedZone("CET", secondsCETOfUTC)

	testCases := []struct {
		Name         string
		Timestamp    string
		IsValid      bool
		ErrorMessage string
		Expected     time.Time
	}{
		{
			"When timestamp is in the future",
			"2488797462",
			true,
			"",
			time.Date(2048, time.November, 12, 13, 37, 42, 0, centralEuropeTime),
		},
		{
			"When timestamp is in the past",
			"942410262",
			true,
			"",
			time.Date(1999, time.November, 12, 13, 37, 42, 0, centralEuropeTime),
		},
		{
			"When timestamp is not a number",
			"12th of November, at 1pm 42",
			false,
			"can't parse unix timestamp",
			time.Now(),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			actual, err := timeutils.FromStringTimestamp(testCase.Timestamp)

			if err == nil && !testCase.IsValid {
				t.Fatalf("Expected error but got none. Received: %+v", actual)
			}

			if err != nil && testCase.IsValid {
				t.Fatalf("Expected no errors but got one. Received: %s", err.Error())
			}

			if err != nil && !strings.Contains(err.Error(), testCase.ErrorMessage) {
				t.Fatalf("Wrong error. Expected an error containing: '%s'. Received: '%s'", testCase.ErrorMessage, err.Error())
			}

			if testCase.IsValid && !actual.Equal(testCase.Expected) {
				t.Fatalf("Wrong date. Expected '%+v'. Received: '%+v'", testCase.Expected, actual)
			}
		})
	}
}
