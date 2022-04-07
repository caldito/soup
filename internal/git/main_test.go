package git

import (
	"testing"
	"fmt"
	"time"
	"sort"
	gogit "github.com/go-git/go-git/v5"
)

func TestGetBranchNamesSoupTest(t *testing.T) {
	// prepare
	cloneLocation = fmt.Sprintf("%s%d", "/tmp/souptest/", time.Now().Unix())
	r, err := gogit.PlainClone(cloneLocation, false, &gogit.CloneOptions{
		URL: "https://github.com/caldito/soup-test.git",
	})
	if err != nil {
		t.Fatalf(`could not clone test repo`)
	}
	// Test
	result := getBranchNames(r)
	if len(result) != 4 {
		t.Fatalf(`expecting 4 branches but returning ` + fmt.Sprint(len(result)))
	}
	sort.Strings(result)
	expectedSlice := []string{"bugs/2", "features/1", "features/2", "main"}
	for i, v := range expectedSlice {
		if v != result[i] {
			t.Fatalf(`expected branches mismatch`)
		}
	}
}