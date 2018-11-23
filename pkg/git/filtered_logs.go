package git

func SquashedCommits(git Git, from Reference, to Reference) ([]*Commit, error) {
	return git.Log(from, to)
}

func MergedCommits(git Git, from Reference, to Reference) ([]*Commit, error) {
	var commits []*Commit

	allCommits, err := git.Log(from, to)
	if err != nil {
		return nil, err
	}

	for _, commit := range allCommits {
		if commit.IsMerge {
			commits = append(commits, commit)
		}
	}

	return commits, nil
}
