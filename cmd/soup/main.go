package main

import (
	"fmt"
	git "github.com/go-git/go-git/v5"
	config "github.com/go-git/go-git/v5/config"
	plumbing "github.com/go-git/go-git/v5/plumbing"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

// Configuration file structs and function

type Namespace struct {
	Namespace string
	Branch    string
}

type Conf struct {
	Namespaces  []Namespace
	Files       []string
	Directories []string
}

func getConf(cloneLocation string) Conf {
	var c Conf
	yamlFile, err := ioutil.ReadFile(cloneLocation + "/.soup.yml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		panic(err)
	}
	return c
}

// Auxiliary functions

func getBranchNames(r *git.Repository) []string {
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
	return branchNames
}

func getNamespace(branchName string, conf Conf) string {
	var namespace string = ""
	for _, a := range conf.Namespaces{
		if (a.Branch == branchName){
			namespace = a.Namespace
		}
	}
	return namespace
}

// Core functions

func deploy(branchName string, cloneLocation string) error {
	// Get configuration from file
	var conf Conf = getConf(cloneLocation)
	var namespace string = getNamespace(branchName, conf)
	if (namespace == ""){
		fmt.Println("Branch " + branchName + " does not match with any namespace to be deployed")
		return nil
	}
	fmt.Println("Deploying branch " + branchName + " to namespace " + namespace)
	return nil
}

func run() error {
	// Clone repo
	cloneLocation := fmt.Sprintf("%s%d", "/tmp/soup/", time.Now().Unix())
	r, err := git.PlainClone(cloneLocation, false, &git.CloneOptions{
		URL: "https://github.com/caldito/soup-test",
	})
	if err != nil {
		panic(err)
	}
	// Get branch names
	branchNames := getBranchNames(r)
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
		err = deploy(branchName, cloneLocation)
		if err != nil {
			panic(err)
		}
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
