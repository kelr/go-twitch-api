package main

import (
	"./twitchapi"
	"encoding/json"
	"fmt"
)

// Provide your Client ID and secret. Set your redirect URI to one that you own.
// Better to set these as environment variables.
const (
	clientID     = ""
	clientSecret = ""
	redirectURI  = "https://twitch.tv"
	targetUsername = ""
)

// Set scopes to request from the user
var scopes = []string{"user:read:email"}

func main() {
	// Setup OAuth2 configs and get the URL to send to the user to ask for perms
	config, url := twitchapi.UserAuthSetup(clientID, clientSecret, redirectURI, &scopes)
	fmt.Println(url)

	// Enter the code received by the redirect URI
	var code string
	if _, err := fmt.Scan(&code); err != nil {
	    fmt.Println(err)
	}

	// User token will be automatically refreshed as long as the client is online.
	client, err := twitchapi.NewTwitchClientUserAuth(config, code)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Get user information, will include email for the user you have a token from
	opt := &twitchapi.GetUsersOpt{
		Login: targetUsername,
	}

	response, err := client.GetUsers(opt)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Pretty print
	obj, _ := json.MarshalIndent(response, "", "  ")
	fmt.Println(string(obj))
}
