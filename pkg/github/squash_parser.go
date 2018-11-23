package github

import (
	"fmt"
	"regexp"

	"github.com/kdisneur/changelog/pkg/parser"
)

type squashParser struct{}

var squashRegex = regexp.MustCompile("\\(#([0-9]+)\\)")

func NewSquashParser() parser.Parser {
	return squashParser{}
}

func (s squashParser) Equal(other parser.Parser) bool {
	_, hasGoodType := other.(squashParser)

	return hasGoodType
}

func (s squashParser) FindID(subject string) (string, error) {
	matches := squashRegex.FindStringSubmatch(subject)

	if len(matches) == 2 {
		return matches[1], nil
	}

	return "", fmt.Errorf("can't parse commit subject '%s'", subject)
}

func (s squashParser) KeepCommit(subject string) bool {
	return squashRegex.MatchString(subject)
}
