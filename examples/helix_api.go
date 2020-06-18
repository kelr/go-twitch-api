package main

import (
	"github.com/kelr/gundyr"
	"log"
)

const (
	clientID     = ""
	clientSecret = ""
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
