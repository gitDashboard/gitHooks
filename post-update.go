package main

import (
	"fmt"
	git "gopkg.in/libgit2/git2go.v22"
	"os"
	"path/filepath"
)

func main() {
	os.Exit(generateGraph())
}

func repoWalker(commit *git.Commit) bool {
	fmt.Printf("parent: %+v, message:%s\n", commit.Parent(0), commit.Summary())
	return true
}

func generateGraph() int {
	gitDir := os.Getenv("GIT_DIR")
	if gitDir == "" {
		fmt.Println("Error GIT_DIR env not found")
		return 1
	}
	gitDir, _ = filepath.Abs(gitDir)
	fmt.Println("gitDir:", gitDir)
	repo, err := git.OpenRepository(gitDir)
	if err != nil {
		fmt.Printf("Error:%s\n", err.Error())
		return 1
	}
	defer repo.Free()
	walk, err := repo.Walk()
	if err != nil {
		fmt.Printf("Error:%s\n", err.Error())
		return 1
	}
	defer walk.Free()

	//finding all branch
	/*refIt, err := repo.NewReferenceIterator()
	if err != nil {
		fmt.Printf("Error:%s\n", err.Error())
		return 1
	}

	defer refIt.Free()
	refNameIt := refIt.Names()
	refName, refNameErr := refNameIt.Next()
	for refNameErr == nil {
		walk.PushRef(refName)
		refName, refNameErr = refNameIt.Next()
	}*/
	walk.PushRef("refs/heads/provaBranch")
	walk.Sorting(git.SortTopological | git.SortTime | git.SortReverse)
	walk.Iterate(repoWalker)
	return 0
}
