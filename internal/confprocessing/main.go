package confprocessing

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"regexp"
	"strings"
)

type Namespace struct {
	Namespace string
	Branch    string
}

type BuildConf struct {
	Namespaces []Namespace
	Manifests  []string
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
		var emptyarr []string
		return "", emptyarr, nil
	}
	// Process configuration
	var namespace string = getNamespace(branchName, buildConf)
	if namespace == "" {
		fmt.Println("Branch " + branchName + " does not match with any namespace to be deployed")
		var emptyarray []string
		return "", emptyarray, nil
	}
	fmt.Println("Deploying branch " + branchName + " to namespace " + namespace)
	return namespace, buildConf.Manifests, nil
}
