// An example to obtain a user's email with a user access token.
package main

import (
	"log"

	"github.com/kelr/gundyr"
	"github.com/kelr/gundyr/auth"
)

// Provide your Client ID and secret. Set your redirect URI to one that you own.
// The URI must match exactly with the one registered by your app on the Twitch Developers site
const (
	clientID     = ""
	clientSecret = ""
	redirectURI  = "http://localhost"
	tokenFile    = "token.json"
)

func main() {
	var scopes = []string{"user:read:email"}
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
	cfg := &gundyr.HelixConfig{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       scopes,
		RedirectURI:  redirectURI,
		Token:        token,
	}
	c, err := gundyr.NewHelix(cfg)
	if err != nil {
		log.Fatal(err)
	}

	email, err := c.GetUserEmail("your-username")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(email)
}
