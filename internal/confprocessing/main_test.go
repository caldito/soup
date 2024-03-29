package confprocessing

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestGetManifests(t *testing.T) {
	// preparing files and directories for test
	cloneLocation := fmt.Sprintf("%s%d", "/tmp/souptestsTestFindFilesByRegex", time.Now().Unix())
	os.Mkdir(cloneLocation, 0755)
	defer os.RemoveAll(cloneLocation)
	f, err := os.Create(cloneLocation + "/a.yml")
	if err != nil {
		t.Fatalf(`failed preparing test`)
		panic(err)
	}
	f.Close()
	f, err = os.Create(cloneLocation + "/b.yml")
	if err != nil {
		t.Fatalf(`failed preparing test`)
		panic(err)
	}
	f.Close()
	os.Mkdir(cloneLocation+"/c", 0755)
	f, err = os.Create(cloneLocation + "/c/d.yml")
	if err != nil {
		t.Fatalf(`failed preparing test`)
		panic(err)
	}
	// test
	manifests, err := getManifests(cloneLocation, []string{"*.yml", "*.yml", "*/*.yml"})
	if err != nil {
		t.Fatalf(`findFilesByRegex returning an error`)
	}
	if len(manifests) != 3 {
		t.Fatalf(`expencting 3 manifests but returning ` + fmt.Sprint(len(manifests)))
	}
	expectedSlice := []string{cloneLocation + "/a.yml", cloneLocation + "/b.yml", cloneLocation + "/c/d.yml"}
	for i, v := range expectedSlice {
		if v != manifests[i] {
			t.Fatalf(`expected manifests mismatch`)
		}
	}
}

func TestRemoveDuplicateStr(t *testing.T) {
	result := removeDuplicateStr([]string{"aa", "aa", "bb", "aa"})
	if len(result) != 2 {
		t.Fatalf(`expencting 2 manifests but returning ` + fmt.Sprint(len(result)))
	}
	expectedSlice := []string{"aa", "bb"}
	for i, v := range expectedSlice {
		if v != result[i] {
			t.Fatalf(`expected slice mismatch`)
		}
	}
}

func TestGetNamespaceExactBranchName(t *testing.T) {
	var ns Namespace;
	ns.Namespace = "prod"
	ns.Branch = "main"
	result, err := getNamespace("main",[]Namespace{ns})
	if err != nil {
		t.Fatalf(`getNamespace returning an error`)
	}
	if result != "prod" {
		t.Fatalf(`Unexpected result: ` + result )
	}
}

func TestGetNamespaceRegex(t *testing.T) {
	var ns Namespace;
	ns.Namespace = "features"
	ns.Branch = "features/*"
	result, err := getNamespace("features/38",[]Namespace{ns})
	if err != nil {
		t.Fatalf(`getNamespace returning an error`)
	}
	if result != "features" {
		t.Fatalf(`Unexpected result: ` + result )
	}
}

func TestGetNamespaceAsBranch(t *testing.T) {
	var ns Namespace;
	ns.Namespace = "as-branch"
	ns.Branch = "features/*"
	result, err := getNamespace("features/38",[]Namespace{ns})
	if err != nil {
		t.Fatalf(`getNamespace returning an error`)
	}
	if result != "features-38" {
		t.Fatalf(`Unexpected result: ` + result )
	}
}

func TestGetBuildConfNoRepoDirectory(t *testing.T) {
	cloneLocation := fmt.Sprintf("%s%d", "/tmp/souptests/TestGetBuildConfNoRepoDirectory", time.Now().Unix())
	_, err := getBuildConf(cloneLocation)
	if err == nil {
		t.Fatalf(`getBuildConf not failing when soup.yml is unexistent in repository`)
	}
}
