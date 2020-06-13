package main

import (
	"fmt"
	"github.com/kelr/go-twitch-api"
)

const (
	clientID     = ""
	clientSecret = ""
)

func main() {
	client, err := helix.NewHelixClient(clientID, clientSecret)
	if err != nil {
		fmt.Println(err)
		return
	}

	userId, err := client.UserToId("kyrotobi")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(userId)

	userName, err := client.IdToUser("31903323")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(userName)
}
