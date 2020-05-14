package helix 

import (
	"fmt"
	"context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/twitch"
)

// Creates and returns OAuth2 configuration object with the twitch endpoint. Also returns a URL to be sent to the user used to initiate authentication.
func NewUserAuth(clientID string, clientSecret string, redirectURI string, scopes *[]string) (*oauth2.Config, string) {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       *scopes,
		Endpoint:     twitch.Endpoint,
		RedirectURL:  redirectURI,
	}
	return config, config.AuthCodeURL("state", oauth2.AccessTypeOffline)
}

// Conducts the exchange to turn an auth code into a user token. The OAuth2 config used to create the auth code must be the same.
func TokenExchange(config *oauth2.Config, authCode string) (*oauth2.Token, error) {
	token, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		fmt.Println("Error in obtaining user token:", err)
		return nil, err
	}
	return token, nil
}
