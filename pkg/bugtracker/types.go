package bugtracker

type BugTracker interface {
	Equal(other BugTracker) bool
	FindIssue(id string) (*Issue, error)
}

type Issue struct {
	ID      string
	Subject string
	Link    string
}

func (i *Issue) Equal(other *Issue) bool {
	return i.ID == other.ID && i.Subject == other.Subject && i.Link == other.Link
}
