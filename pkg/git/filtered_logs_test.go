package git_test

import (
	"testing"
	"time"

	"github.com/kdisneur/changelog/pkg/git"
	"github.com/kdisneur/changelog/pkg/testing/repository"
)

func TestMergedCommits(t *testing.T) {
	gitRepository := repository.New("git@github.com:kdisneur/changelog")

	john := git.Person{Fullname: "John Doe", Email: "john.doe@gmail.com"}
	date := time.Now()

	gitRepository.AddCommit("a2bc4fd34ba164ad0c1a264340ce37b0dbdaa6ef", john, date, "Commit 1")
	gitRepository.AddMergeCommit("6398b4e189b94ce300641431d3dfa00c373d1bb1", john, date.Add(time.Hour*2), "Commit 2")
	gitRepository.AddCommit("4f28c412c51c44c94daa3fced544567c3f94dd7b", john, date.Add(time.Hour*3), "Commit 3")
	gitRepository.AddMergeCommit("555475c1e0c506eaf23d0db155f6592f7383c495", john, date.Add(time.Hour*4), "Commit 4")
	gitRepository.AddCommit("16dd9970c4f776157ccc6a7d8c78b2bdeeaab1c4", john, date.Add(time.Hour*5), "Commit 5")

	commits, _ := git.MergedCommits(gitRepository, git.Reference("a2bc4fd34ba164ad0c1a264340ce37b0dbdaa6ef"), git.Reference("16dd9970c4f776157ccc6a7d8c78b2bdeeaab1c4"))

	if len(commits) != 2 {
		t.Fatalf("We should have received 1 commit, received: %+v", commits)
	}

	if commits[0].Message != "Commit 2" {
		t.Errorf("We should have received commit 2 but we received: %+v", commits[0])
	}

	if commits[1].Message != "Commit 4" {
		t.Errorf("We should have received commit 4 but we received: %+v", commits[0])
	}
}

func TestSquashedCommits(t *testing.T) {
	gitRepository := repository.New("git@github.com:kdisneur/changelog")

	john := git.Person{Fullname: "John Doe", Email: "john.doe@gmail.com"}
	date := time.Now()

	gitRepository.AddCommit("a2bc4fd34ba164ad0c1a264340ce37b0dbdaa6ef", john, date, "Commit 1")
	gitRepository.AddMergeCommit("6398b4e189b94ce300641431d3dfa00c373d1bb1", john, date.Add(time.Hour*2), "Commit 2")
	gitRepository.AddCommit("4f28c412c51c44c94daa3fced544567c3f94dd7b", john, date.Add(time.Hour*3), "Commit 3")
	gitRepository.AddMergeCommit("555475c1e0c506eaf23d0db155f6592f7383c495", john, date.Add(time.Hour*4), "Commit 4")
	gitRepository.AddCommit("16dd9970c4f776157ccc6a7d8c78b2bdeeaab1c4", john, date.Add(time.Hour*5), "Commit 5")

	commits, _ := git.SquashedCommits(gitRepository, git.Reference("a2bc4fd34ba164ad0c1a264340ce37b0dbdaa6ef"), git.Reference("16dd9970c4f776157ccc6a7d8c78b2bdeeaab1c4"))

	expectedNumberCommits := len(gitRepository.Commits) - 1
	if expectedNumberCommits != len(commits) {
		t.Fatalf("We should have received %d commits, received: %d\nExpected:\n%+v\n\nReceived:\n%+v", expectedNumberCommits, len(commits), commits, gitRepository.Commits)
	}

	for index, expectedCommit := range commits {
		actualCommit := commits[index]

		if !expectedCommit.Equal(actualCommit) {
			t.Errorf("Wrong commit.\nExpected:\n%+v\nReceived:\n%+v", expectedCommit, actualCommit)
		}
	}
}
