package main

import (
	"github.com/kelr/gundyr"
	"log"
)

const (
	clientID     = "v1jznhyjrk89g65v6if0jpymwk7s4e"
	clientSecret = "qlf6iyrg33xsxcx0l5khkgqfecf7a0"
)

func main() {
	c, err := gundyr.NewHelix(clientID, clientSecret)
	if err != nil {
		log.Fatal(err)
	}

	userID, err := c.UserToID("kyrotobi")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(userID)

	userName, err := c.IDToUser(userID)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(userName)
}
