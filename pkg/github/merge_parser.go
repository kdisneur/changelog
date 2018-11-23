package github

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/kdisneur/changelog/pkg/parser"
)

type mergeParser struct{}

var mergeRegex = regexp.MustCompile("Merge pull request #([0-9]+)")

func NewMergeParser() parser.Parser {
	return mergeParser{}
}

func (m mergeParser) FindID(subject string) (string, error) {
	matches := mergeRegex.FindStringSubmatch(subject)

	if len(matches) == 2 {
		return matches[1], nil
	} else {
		return "", errors.New(fmt.Sprintf("can't parse merge subject '%s'", subject))
	}
}

func (m mergeParser) KeepCommit(subject string) bool {
	return mergeRegex.MatchString(subject)
}

func (m mergeParser) Equal(other parser.Parser) bool {
	_, hasGoodType := other.(mergeParser)

	return hasGoodType
}
