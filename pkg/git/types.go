package git

import (
	"fmt"
	"time"
)

type Path string

func (p Path) String() string {
	return string(p)
}

type Reference string

type Person struct {
	Fullname string
	Email    string
}

func (p Person) Equal(other Person) bool {
	return p.Fullname == other.Fullname && p.Email == other.Email
}

type Commit struct {
	ID          string
	Author      Person
	AuthoredAt  time.Time
	Committer   Person
	CommittedAt time.Time
	IsMerge     bool
	Message     string
}

func (c *Commit) String() string {
	return fmt.Sprintf(
		`commit %s (merge %t)
Author: %s <%s>
Date: %s

Commiter: %s <%s>
Date: %s

%s`,
		c.ID, c.IsMerge,
		c.Author.Fullname, c.Author.Email,
		c.AuthoredAt,
		c.Committer.Fullname, c.Committer.Email,
		c.CommittedAt,
		c.Message)
}

func (c *Commit) Equal(other *Commit) bool {
	return c.ID == other.ID &&
		c.Author.Equal(other.Author) &&
		c.AuthoredAt.Equal(other.AuthoredAt) &&
		c.Committer.Equal(other.Committer) &&
		c.CommittedAt.Equal(other.CommittedAt) &&
		c.IsMerge == other.IsMerge &&
		c.Message == other.Message
}

type RemoteType int

const (
	HTTPS RemoteType = iota
	GIT
)

type Remote struct {
	Type           RemoteType
	Host           string
	RepositoryName string
}

type Git interface {
	Equal(other Git) bool
	Log(from Reference, to Reference) ([]*Commit, error)
	FindRemote() (*Remote, error)
}
