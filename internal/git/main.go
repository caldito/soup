package git

import (
	// from this repository
	"github.com/caldito/soup/internal/confprocessing"
	"github.com/caldito/soup/internal/deployment"

	// from other repositories
	"fmt"
	gogit "github.com/go-git/go-git/v5"
	config "github.com/go-git/go-git/v5/config"
	plumbing "github.com/go-git/go-git/v5/plumbing"
	"os"
	"strings"
	"time"
)

var cloneLocation string

// Auxiliary functions
func getBranchNames(r *gogit.Repository) []string {
	var branchNames []string
	remote, err := r.Remote("origin")
	if err != nil {
		fmt.Println("Error getting remote origin")
		panic(err)
	}
	refList, err := remote.List(&gogit.ListOptions{})
	if err != nil {
		fmt.Println("Error getting branch list")
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
	return branchNames
}

// Core functions

func processBranch(branchName string) error {
	// Deployment configuration processing
	namespace, manifests, err := confprocessing.ProcessConf(branchName, cloneLocation)
	if err != nil {
		fmt.Println("Error processing branch")
		panic(err)
	}
	if namespace == "" {
		return nil
	}
	// Deployment module
	err = deployment.Deploy(namespace, manifests)
	if err != nil {
		fmt.Println("Error deploying")
		panic(err)
	}
	return nil
}

func LoopBranches(repo string) error {
	// Clone repo
	cloneLocation = fmt.Sprintf("%s%d", "/tmp/soup/", time.Now().Unix())
	r, err := gogit.PlainClone(cloneLocation, false, &gogit.CloneOptions{
		URL: repo,
	})
	if err != nil {
		fmt.Println("Error downloading repo")
		panic(err)
	}
	// Get branch names
	branchNames := getBranchNames(r)
	// Fetch branches
	err = r.Fetch(&gogit.FetchOptions{
		RefSpecs: []config.RefSpec{"refs/*:refs/*", "HEAD:refs/heads/HEAD"},
	})
	if err != nil {
		fmt.Println("Error fetching branches")
		panic(err)
	}
	// Checkout to the branches and do GitOps stuff
	w, _ := r.Worktree()
	for _, branchName := range branchNames {
		// Checkout
		err = w.Checkout(&gogit.CheckoutOptions{
			Branch: plumbing.ReferenceName("refs/heads/" + branchName),
			Force:  true,
		})
		if err != nil {
			fmt.Println("Error checking out to " + branchName)
			panic(err)
		}
		// Process branch
		err = processBranch(branchName)
		if err != nil {
			fmt.Println("Error processing branch")
			panic(err)
		}
	}
	os.RemoveAll(cloneLocation)
	return nil
}
