package main

import "fmt"
import "os"
import "github.com/go-git/go-git/v5"
import "github.com/go-git/go-git/v5/storage/memory"
import "github.com/go-git/go-git/v5/plumbing/object"

func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

func main() {
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: "https://github.com/go-git/go-billy",
	})

	CheckIfError(err)

	// ... retrieves the branch pointed by HEAD
	ref, err := r.Head()
	CheckIfError(err)


	// ... retrieves the commit history
	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	CheckIfError(err)

	// ... just iterates over the commits, printing it
	err = cIter.ForEach(func(c *object.Commit) error {
		fmt.Println(c)
		return nil
	})
	CheckIfError(err)

}
