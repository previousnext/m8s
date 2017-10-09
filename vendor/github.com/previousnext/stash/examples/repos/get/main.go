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

	project = flag.String("g", "", "project to get")
)

func main() {
	flag.Parse()
	c := stash.NewClient(*username, *password, *host)
	project, err := c.RepoService().GetRepos(*project)
	if err != nil {
		fmt.Println("failed to retrieve user: ", err)
		return
	}

	pretty.Println(project)
}
