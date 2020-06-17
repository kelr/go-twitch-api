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
		return
	}

	userID, err := c.UserToId("kyrotobi")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(userID)

	userName, err := c.IdToUser("31903323")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(userName)
}
