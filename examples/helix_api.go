package main

import (
	"log"

	"github.com/kelr/gundyr"
)

const (
	clientID     = ""
	clientSecret = ""
)

func main() {
	cfg := &gundyr.HelixConfig{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}

	c, err := gundyr.NewHelix(cfg)
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
