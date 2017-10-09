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

	group = flag.String("g", "", "group to retrieve users for")
)

func main() {
	flag.Parse()
	c := stash.NewClient(*username, *password, *host)
	users, err := c.GroupService().GetUsers(*group, "")
	if err != nil {
		fmt.Println("failed to retrieve groups users: ", err)
		return
	}

	pretty.Println(users)
}
