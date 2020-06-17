// An example to obtain a user's email with a user access token.
package main

import (
	"encoding/json"
	"github.com/kelr/gundyr/auth"
	"github.com/kelr/gundyr/helix"
	"log"
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

// Load an OAuth2 token from a file and check for validity.
// The token must have been created with a matching scope.
func retrieveToken() *oauth2.Token {
	token, err = auth.LoadTokenFile(config, tokenFile)

	if err != nil {
		log.Fatal(err)
	}

	// Verify that the cached token has not expired.
	newToken := auth.VerifyToken(config, token)

	// Update the token file.
	if err := auth.FlushTokenFile(tokenFile, newToken); err != nil {
		log.Fatal(err)
	}
}

func main() {
	scopes := []string{"user:read:email"}

	// Setup OAuth2 configuration
	config, err := auth.NewUserAuth(clientID, clientSecret, redirectURI, &scopes)
	if err != nil {
		log.Fatal(err)
	}

	// Create the API client. User token will be automatically refreshed.
	// See examples/auth_token.go for an example on creating a new token.
	client, err := helix.NewClientUserAuth(config, retrieveToken())
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
