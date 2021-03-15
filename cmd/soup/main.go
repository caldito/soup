package main

import (
	"fmt"
	git "github.com/go-git/go-git/v5"
	config "github.com/go-git/go-git/v5/config"
	plumbing "github.com/go-git/go-git/v5/plumbing"
	"os"
	"strings"
	"time"
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

func deploy(branch string) error {

	return nil
}

func run() error {
	// Clone repo
	cloneLocation := "/tmp/soup/" + string(time.Now().Unix())

	r, err := git.PlainClone(cloneLocation, false, &git.CloneOptions{
		URL: "https://github.com/caldito/soup-test",
	})
	if err != nil {
		panic(err)
	}
	// Get branch names
	branchNames, err := getBranchNames(r)
	if err != nil {
		panic(err)
	}
	// Fetch branches
	err = r.Fetch(&git.FetchOptions{
		RefSpecs: []config.RefSpec{"refs/*:refs/*", "HEAD:refs/heads/HEAD"},
	})
	if err != nil {
		panic(err)
	}
	// Checkout to the branches and do GitOps stuff
	w, _ := r.Worktree()
	for _, branchName := range branchNames {
		err = w.Checkout(&git.CheckoutOptions{
			Branch: plumbing.ReferenceName("refs/heads/" + branchName),
			Force:  true,
		})
		if err != nil {
			panic(err)
		}
		// Deploy after checking branch
		fmt.Println("Deploying branch " + branchName + "...")
		err = deploy(branchName)
		if err != nil {
			panic(err)
		}
		fmt.Println("Branch " + branchName + " deployed")
	}
	os.RemoveAll(cloneLocation)
	return nil
}

func main() {
	for {
		err := run()
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second * 30)
	}
}
