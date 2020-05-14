// An example to obtain a user authentication token for user's email.
// Uses the token to get info about the user.
package main

import (
	"../helix"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"time"
)

// Provide your Client ID and secret. Set your redirect URI to one that you own.
// Better to set these as environment variables.
const (
	clientID     = ""
	clientSecret = ""
	redirectURI  = ""
)

// Set scopes to request from the user
var scopes = []string{"user:read:broadcast"}

func main() {
	// Setup OAuth2 configs and get the URL to send to the user to ask for perms
	config, url := helix.NewUserAuth(clientID, clientSecret, redirectURI, &scopes)
	fmt.Println(url)

	// Import an existing token to use
	token := new(oauth2.Token)
	token.AccessToken = ""
	token.Expiry = time.Date(2020, 5, 14, 6, 45, 0, 0, time.UTC)
	token.RefreshToken = ""
	token.TokenType = "bearer"

	// User token will be automatically refreshed as long as the client is online.
	client, err := helix.NewTwitchClientUserAuth(config, token)
	if err != nil {
		return
	}

	// Get a list of all active extensions for the user matching the token
	resp, err := client.GetUserActiveExtensions(nil)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Pretty print
	obj, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Println(string(obj))
}
