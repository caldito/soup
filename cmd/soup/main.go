package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	git "github.com/go-git/go-git/v5"
	config "github.com/go-git/go-git/v5/config"
	plumbing "github.com/go-git/go-git/v5/plumbing"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	// imports for doSSA
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	yamlk8s "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
)

// Global variables
var programConf ProgramConf

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
			return namespace
		}
	}
	return ""
}

func getBuildConf(cloneLocation string) BuildConf {
	var buildConf BuildConf
	yamlFile, err := ioutil.ReadFile(cloneLocation + "/.soup.yml")
	if err != nil {
		fmt.Println("Error reading .soup file")
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, &buildConf)
	if err != nil {
		fmt.Println("Error unmarshalling build conf")
		panic(err)
	}
	return buildConf
}

func deploy(namespace string, manifests []string) error {
	for _, manifest := range manifests {
		config, err := rest.InClusterConfig()
		if err != nil {
			fmt.Println("Error getting cluster config")
			panic(err)
		}
		ctx := context.TODO()
		err = doSSA(ctx, config, namespace, manifest)
		if err != nil {
			fmt.Println("Error deploying the manifest" + manifest)
			panic(err)
		}
	}
	return nil
}

// TODO export this function to package
func doSSA(ctx context.Context, cfg *rest.Config, namespace string, manifest string) error {
	var decUnstructured = yamlk8s.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)

	// 1. Prepare a RESTMapper to find GVR
	dc, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		return err
	}
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))

	// 2. Prepare the dynamic client
	dyn, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return err
	}

	// 3. Decode YAML manifest into unstructured.Unstructured
	yamlFile, err := ioutil.ReadFile(manifest)
	if err != nil {
		fmt.Println("Error reading manifest" + manifest)
		return err
	}

	obj := &unstructured.Unstructured{}
	_, gvk, err := decUnstructured.Decode(yamlFile, nil, obj)
	if err != nil {
		return err
	}

	// 4. Find GVR
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return err
	}

	// 5. Obtain REST interface for the GVR
	var dr dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		// namespaced resources should specify the namespace
		dr = dyn.Resource(mapping.Resource).Namespace(obj.GetNamespace())
	} else {
		// for cluster-wide resources
		dr = dyn.Resource(mapping.Resource)
	}

	// 6. Marshal object into JSON
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	// 7. Create or Update the object with SSA
	//     types.ApplyPatchType indicates SSA.
	//     FieldManager specifies the field owner ID.
	_, err = dr.Patch(ctx, obj.GetName(), types.ApplyPatchType, data, metav1.PatchOptions{
		FieldManager: "sample-controller",
	})

	return err
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

func processBranch(branchName string, cloneLocation string) error {
	// Get configuration from file
	var buildConf BuildConf = getBuildConf(cloneLocation)
	// Process configuration
	var namespace string = getNamespace(branchName, buildConf)
	if namespace == "" {
		fmt.Println("Branch " + branchName + " does not match with any namespace to be deployed")
		return nil
	}
	fmt.Println("Deploying branch " + branchName + " to namespace " + namespace)
	// Deploy
	err := deploy(namespace, buildConf.Manifests)
	if err != nil {
		fmt.Println("Error deploying")
		panic(err)
	}
	return nil
}

func run() error {
	// Clone repo
	cloneLocation := fmt.Sprintf("%s%d", "/tmp/soup/", time.Now().Unix())
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
		err = processBranch(branchName, cloneLocation)
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
