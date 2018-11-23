package parser

type Parser interface {
	FindID(subject string) (string, error)
	KeepCommit(subject string) bool
	Equal(other Parser) bool
}
