package main

import (
	"os"
	"fmt"
	"strings"
	"time"
	git "github.com/go-git/go-git/v5"
	config "github.com/go-git/go-git/v5/config"
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

func run() error {
	// Clone repo
	cloneLocation := "/tmp/soup"
	os.RemoveAll(cloneLocation)
	r, err := git.PlainClone(cloneLocation, false, &git.CloneOptions{
		URL: "https://github.com/caldito/ipwarn",
	})
	if err != nil {
		panic(err)
	}
	// Get branch names
	branchNames, err := getBranchNames(r)
	if err != nil {
		panic(err)
	}
	fmt.Print(branchNames)
	// Fetch branches
	err = r.Fetch(&git.FetchOptions{
        RefSpecs: []config.RefSpec{"refs/*:refs/*", "HEAD:refs/heads/HEAD"},
    })
    if err != nil {
        panic(err)
    }
	// Checkout to the branches and do GitOps stuff
	w, _ := r.Worktree()
	for _, branch := range branchNames {
		err = w.Checkout(&git.CheckoutOptions{
			Branch: fmt.Sprintf("refs/heads/%s", branch),
			Force: true,
		})
		if err != nil {
			panic(err)
		}
		// TODO GitOps stuff after checking branch
		fmt.Sprintf("checkout to %s", branchName)
	}
	return nil
}

func main() {
	for {
		err := run()
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second*30)
	}
}
