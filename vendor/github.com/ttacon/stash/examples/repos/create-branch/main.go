package main

import (
	"flag"
	"fmt"

	"github.com/ttacon/pretty"
	"github.com/ttacon/stash"
)

var (
	username = flag.String("u", "", "username")
	password = flag.String("p", "", "password")
	host     = flag.String("h", "http://localhost:7990", "host to query")

	project    = flag.String("g", "", "project to get")
	repoName   = flag.String("r", "", "name of repo to create")
	branchName = flag.String("b", "", "branch to create")
)

func main() {
	flag.Parse()
	c := stash.NewClient(*username, *password, *host)
	project, err := c.RepoService().CreateBranch(
		*project, *repoName, *branchName, "refs/heads/master")
	if err != nil {
		fmt.Println("failed to retrieve user: ", err)
		return
	}

	pretty.Println(project)
}
