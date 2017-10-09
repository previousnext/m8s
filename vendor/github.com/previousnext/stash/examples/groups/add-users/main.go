package main

import (
	"flag"
	"fmt"

	"github.com/ttacon/stash"
)

var (
	username = flag.String("u", "", "username")
	password = flag.String("p", "", "password")
	host     = flag.String("h", "http://localhost:7990", "host to query")

	group     = flag.String("g", "", "group to update")
	userToAdd = flag.String("u2a", "", "user to add to group")
)

func main() {
	flag.Parse()
	c := stash.NewClient(*username, *password, *host)
	err := c.GroupService().AddUsers(*group, *userToAdd)
	if err != nil {
		fmt.Println("failed to add user to group: ", err)
		return
	}

	fmt.Println("successfully added user")
}
