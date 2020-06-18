// An example to obtain a user's email with a user access token.
package main

import (
	"encoding/json"
	"github.com/kelr/gundyr/auth"
	"github.com/kelr/gundyr/helix"
	"log"
	"fmt"
)

// Provide your Client ID and secret. Set your redirect URI to one that you own.
// The URI must match exactly with the one registered by your app on the Twitch Developers site
const (
	clientID       = ""
	clientSecret   = ""
	redirectURI    = "http://localhost"
	targetUsername = ""
	tokenFile      = "token.json"
)

func main() {
	scopes := []string{"user:read:email"}

	// Setup OAuth2 configuration
	config, err := auth.NewUserAuth(clientID, clientSecret, redirectURI, &scopes)
	if err != nil {
		log.Fatal(err)
	}
	
	// See examples/auth_token.go for an example on creating a new token.
	token, err := auth.RetrieveTokenFile(config, tokenFile)
	if err != nil {
		log.Fatal(err)
	}
	
	// Create the API client. User token will be automatically refreshed.
	client, err := helix.NewClientUserAuth(config, token)
	if err != nil {
		log.Fatal(err)
	}

	// Get user information, will include email for the user you have a token from.
	response, err := client.GetUsers(&helix.GetUsersOpt{Login: targetUsername})
	if err != nil {
		log.Fatal(err)
	}

	// Pretty print
	obj, _ := json.MarshalIndent(response, "", "  ")
	fmt.Println(string(obj))
}
