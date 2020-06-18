// Example to show how to generate a new OAuth2 token using the 3-legged flow.
// The user of this package is responsible for the redirectURI and the transmission
// of the Auth Code URL to the end user.
package main

import (
	"flag"
	"fmt"
	"github.com/kelr/gundyr/auth"
	"golang.org/x/oauth2"
	"log"
)

const (
	clientID     = ""
	clientSecret = ""
	redirectURI  = "http://localhost"
	tokenFile    = "token.json"
)

// Set scopes to request from the user
var scopes = []string{"channel:read:redemptions"}

// Run with -a to generate a new user token
var doAuth = flag.Bool("a", false, "Generate a URL for user token authentication")

// Helper function to generate the auth code URL, generate a user credential token and flush it to a file.
func generateNewToken(config *oauth2.Config) (*oauth2.Token, error) {
	// Get the URL to send to the user and the state code to protect against CSRF attacks.
	url, state := auth.GetAuthCodeURL(config)
	fmt.Println(url)
	fmt.Println("Ensure that state recieved at URI is:", state)

	// Enter the code received by the redirect URI. Ensure that the state value
	// obtained at the redirect URI matches the previous state value.
	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, err
	}

	// Obtain the user token through the code. This token can be reused as long as
	// it has not expired, but the auth code cannot be reused.
	token, err := auth.TokenExchange(config, authCode)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func main() {
	flag.Parse()

	// Setup OAuth2 configuration
	config, err := auth.NewUserAuth(clientID, clientSecret, redirectURI, &scopes)
	if err != nil {
		log.Fatal(err)
	}

	var token *oauth2.Token
	// Generate a new token if the -a flag was provided
	if *doAuth {
		token, err = generateNewToken(config)
	} else {
		token, err = auth.LoadTokenFile(tokenFile)
	}
	if err != nil {
		log.Fatal(err)
	}

	// Verify that the cached token has not expired. Refresh if so.
	newToken := auth.VerifyToken(config, token)
	if newToken.AccessToken != token.AccessToken {
		fmt.Println("Token was expired, now refreshed.")
	}

	// Write the token to a file
	if err := auth.FlushTokenFile(tokenFile, newToken); err != nil {
		log.Fatal(err)
	}

	fmt.Println(newToken.AccessToken)
	fmt.Println(newToken.TokenType)
	fmt.Println(newToken.RefreshToken)
	fmt.Println(newToken.Expiry)
}
