package main

import "fmt"
import "github.com/go-git/go-git/v5"
import "github.com/go-git/go-git/v5/storage/memory"
import "github.com/go-git/go-git/v5/plumbing"
import "github.com/go-git/go-git/v5/plumbing/storer"


func getRemoteBranches(s storer.ReferenceStorer) (storer.ReferenceIter, error) {
	refs, err := s.IterReferences()
	if err != nil {
		return nil, err
	}

	return storer.NewReferenceFilteredIter(func(ref *plumbing.Reference) bool {
		return ref.Name().IsRemote()
	}, refs), nil
}

func main() {
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: "https://github.com/caldito/ipwarn",
	})
	if err != nil {
		panic(err)
	}

	branches, err := getRemoteBranches(r.Storer)
	if err != nil {
		panic(err)
	}

	err = branches.ForEach(func(b *plumbing.Reference) error {
		fmt.Println(b)
		return nil
	})
	if err != nil {
		panic(err)
	}
}
