package confprocessing

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

type Namespace struct {
	Namespace string
	Branch    string // regex
}

type BuildConf struct {
	Namespaces []Namespace
	Manifests  []string // pattern
}

func getNamespace(branchName string, namespaces []Namespace) (string, error) {
	var namespace string = ""
	for _, a := range namespaces {
		matched, err := regexp.MatchString(a.Branch, branchName)
		if err != nil {
			fmt.Println("Error matching strings to get namespace. Branch regex \"" + a.Branch + "\" not valid")
			continue
		}
		if matched {
			if a.Namespace == "as-branch" {
				namespace = branchName
			} else {
				namespace = a.Namespace
			}
			namespace = strings.ReplaceAll(namespace, "/", "-")
			return namespace, nil
		}
	}
	return "", nil
}

func getManifests(cloneLocation string, manifestPatterns []string) ([]string, error) {
	var processedManifests []string
	for _, manifestPattern := range manifestPatterns {
		manifests, err := filepath.Glob(cloneLocation + "/" + manifestPattern)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Error processing the following manifests regex:" + manifestPattern +
				"  You can find how to form the regex here: https://pkg.go.dev/path/filepath#Match")
		} else {
			processedManifests = append(processedManifests, manifests...)
		}
	}
	processedManifests = removeDuplicateStr(processedManifests)
	return processedManifests, nil
}

func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func getBuildConf(cloneLocation string) (BuildConf, error) {
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

func ProcessConf(branchName string, cloneLocation string) (string, []string, error) {
	var buildConf BuildConf
	buildConf, err := getBuildConf(cloneLocation)
	if err != nil {
		// print no build conf found
		fmt.Println("Skipping branch " + branchName + ": Error reading or parsing file .soup.yml")
		return "", nil, nil
	}
	namespace, _ := getNamespace(branchName, buildConf.Namespaces)
	if namespace == "" {
		fmt.Println("Branch " + branchName + " does not match with any namespace to be deployed")
		return "", nil, nil
	}
	processedManifests, err := getManifests(cloneLocation, buildConf.Manifests)
	if err != nil {
		fmt.Println("Skipping branch " + branchName + ": Error getting manifest files")
		return "", nil, nil
	}
	fmt.Println("Deploying branch " + branchName + " to namespace " + namespace)
	return namespace, processedManifests, nil
}
