package main

import (
	"fmt"
	"strings"
	git "github.com/go-git/go-git/v5"
)


func getBranchNames(r *git.Repository) ([]string, error) {
	var branchNames []string
	remote, err := r.Remote("origin")
	if err != nil {
		panic(err)
	}
	refList, err := remote.List(&git.ListOptions{})
	if err != nil {
		panic(err)
	}
	refPrefix := "refs/heads/"
	for _, ref := range refList {
		refName := ref.Name().String()
		if !strings.HasPrefix(refName, refPrefix) {
			continue
		}
		branchName := refName[len(refPrefix):]
		branchNames = append(branchNames, branchName)
	}
	return branchNames, nil
}

func main() {
	// Clone repo
	cloneLocation := "/tmp/soup"
	r, err := git.PlainClone(cloneLocation, false, &git.CloneOptions{
		URL: "https://github.com/caldito/ipwarn",
	})
	if err != nil {
		panic(err)
	}

	branchNames, err := getBranchNames(r)
	if err != nil {
		panic(err)
	}
	fmt.Print(branchNames)
}
