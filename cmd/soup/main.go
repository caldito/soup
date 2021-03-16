package main

import (
	"fmt"
	git "github.com/go-git/go-git/v5"
	config "github.com/go-git/go-git/v5/config"
	plumbing "github.com/go-git/go-git/v5/plumbing"
	"os"
	"strings"
	"time"
	"gopkg.in/yaml.v2"
    "io/ioutil"
)

type namespace struct {
	namespace string `yaml:"namespace"`
	branch string `yaml:"branch"`
}

type conf struct {
    namespaces []namespace `yaml:"namespaces"`
    files []string `yaml:"files"`
	directories []string `yaml:"directories"`
}

func (c *conf) getConf(cloneLocation string) *conf {
    yamlFile, err := ioutil.ReadFile(cloneLocation + "/.soup.yml")
    if err != nil {
        panic(err)
    }
    err = yaml.Unmarshal(yamlFile, c)
    if err != nil {
        panic(err)
    }

    return c
}

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

func deploy(branch string, cloneLocation string) error {
    var c conf
    c.getConf(cloneLocation)

    fmt.Println(c)
	return nil
}

func run() error {
	// Clone repo
	cloneLocation := fmt.Sprintf("%s%d", "/tmp/soup/", time.Now().Unix())
	//cloneLocation := string("/" + time.Now().Unix())
	fmt.Println(cloneLocation)
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
		err = deploy(branchName, cloneLocation)
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
