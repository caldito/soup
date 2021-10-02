package main

import (
	// from this repository
	"github.com/caldito/soup/internal/git"

	// from other places
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	var repo string
	var interval int
	flag.StringVar(&repo, "repo", "", "url of the repository")
	flag.IntVar(&interval, "interval", 120, "execution interval")
	flag.Parse()
	if repo == "" {
		fmt.Println("Exiting, repo flag is not provided")
		os.Exit(2)
	}
	for {
		err := git.LoopBranches(repo)
		if err != nil {
			fmt.Println("Error in run() method")
			os.Exit(1)
		}
		fmt.Printf("%s%d%s", "Sleeping ", interval, "s until next execution...\n\n")
		time.Sleep(time.Second * time.Duration(interval))
	}
}
