package configuration

func (f File) FindRepository(name string) (*GitRepository, bool) {
	for _, repository := range f.Repository {
		if repository.Name == name {
			return &repository, true
		}
	}

	return nil, false
}
