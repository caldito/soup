package main

import (
	// from this repository
	"github.com/caldito/soup/pkg/k8s"

	// from other places
	"context"
	"flag"
	"fmt"
	git "github.com/go-git/go-git/v5"
	config "github.com/go-git/go-git/v5/config"
	plumbing "github.com/go-git/go-git/v5/plumbing"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"k8s.io/client-go/rest"
	"os"
	"regexp"
	"strings"
	"time"
)

// Global variables
var programConf ProgramConf
var cloneLocation string

// Structs
type Namespace struct {
	Namespace string
	Branch    string
}

type BuildConf struct {
	Namespaces []Namespace
	Manifests  []string
}

type ProgramConf struct {
	Repo     string
	Interval int
}

// Auxiliary functions
func getBranchNames(r *git.Repository) []string {
	var branchNames []string
	remote, err := r.Remote("origin")
	if err != nil {
		fmt.Println("Error getting remote origin")
		panic(err)
	}
	refList, err := remote.List(&git.ListOptions{})
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

func getNamespace(branchName string, buildConf BuildConf) string {
	var namespace string = ""
	for _, a := range buildConf.Namespaces {
		matched, err := regexp.MatchString(a.Branch, branchName)
		if err != nil {
			fmt.Println("Error matching strings to get namespace")
			panic(err)
		}
		if matched {
			if a.Namespace == "as-branch" {
				namespace = branchName
			} else {
				namespace = a.Namespace
			}
			namespace = strings.ReplaceAll(namespace, "/", "-")
			return namespace
		}
	}
	return ""
}

func getBuildConf() (BuildConf, error) {
	var buildConf BuildConf
	yamlFile, err := ioutil.ReadFile(cloneLocation + "/.soup.yml")
	if err != nil {
		return buildConf, err
	}
	err = yaml.Unmarshal(yamlFile, &buildConf)
	if err != nil {
		return buildConf, err
	}
	return buildConf, err
}

func deploy(namespace string, manifests []string) error {
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Println("Error getting cluster config")
		panic(err)
	}
	ctx := context.TODO()
	err = k8s.DeclareNamespaceSSA(ctx, config, namespace)
	if err != nil {
		fmt.Println("Error preparing namespace " + namespace)
		panic(err)
	}
	for _, manifest := range manifests {
		err = k8s.DoSSA(ctx, config, namespace, cloneLocation+"/"+manifest)
		if err != nil {
			fmt.Println("Error deploying the manifest " + manifest)
			panic(err)
		}
	}
	return nil
}

// Core functions
func init() {
	flag.StringVar(&programConf.Repo, "repo", "", "url of the repository")
	flag.IntVar(&programConf.Interval, "interval", 120, "execution interval")
	flag.Parse()
	if programConf.Repo == "" {
		fmt.Println("Exiting, repo flag is not provided")
		os.Exit(1)
	}
}

func processBranch(branchName string) error {
	// Get configuration from file
	var buildConf BuildConf 
	buildConf, err := getBuildConf()
	if err != nil {
		// print no build conf found
		fmt.Println("Skipping branch "+ branchName + ": Error reading or parsing file .soup.yml" )
		return nil
	}
	// Process configuration
	var namespace string = getNamespace(branchName, buildConf)
	if namespace == "" {
		fmt.Println("Branch " + branchName + " does not match with any namespace to be deployed")
		return nil
	}
	fmt.Println("Deploying branch " + branchName + " to namespace " + namespace)
	// Deploy
	err = deploy(namespace, buildConf.Manifests)
	if err != nil {
		fmt.Println("Error deploying")
		panic(err)
	}
	return nil
}

func run() error {
	// Clone repo
	cloneLocation = fmt.Sprintf("%s%d", "/tmp/soup/", time.Now().Unix())
	r, err := git.PlainClone(cloneLocation, false, &git.CloneOptions{
		URL: programConf.Repo,
	})
	if err != nil {
		fmt.Println("Error downloading repo")
		panic(err)
	}
	// Get branch names
	branchNames := getBranchNames(r)
	// Fetch branches
	err = r.Fetch(&git.FetchOptions{
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
		err = w.Checkout(&git.CheckoutOptions{
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
	fmt.Printf("%s%d%s", "Sleeping ", programConf.Interval, "s until next execution...\n\n")
	time.Sleep(time.Second * time.Duration(programConf.Interval))
	return nil
}

func main() {
	for {
		err := run()
		if err != nil {
			fmt.Println("Error in run() method")
			panic(err)
		}
	}
}
