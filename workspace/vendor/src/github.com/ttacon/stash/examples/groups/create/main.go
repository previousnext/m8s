package main

import (
	"flag"

	"github.com/ttacon/stash"
)

var (
	username = flag.String("u", "", "username")
	password = flag.String("p", "", "password")
	host     = flag.String("h", "http://localhost:7990", "host to query")

	group = flag.String("g", "", "group to create")
)

func main() {
	flag.Parse()
	c := stash.NewClient(*username, *password, *host)
	c.GroupService().CreateGroup(*group)
}
