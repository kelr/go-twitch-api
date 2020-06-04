package helix

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/twitch"
)

const (
   stateSet = "abcdefghijklmnopqrstuvwxyz123456789"
)
// Creates and returns OAuth2 configuration object with the twitch endpoint. Also returns a URL to be sent to the user used to initiate authentication.
func NewUserAuth(clientID string, clientSecret string, redirectURI string, scopes *[]string) (*oauth2.Config) {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       *scopes,
		Endpoint:     twitch.Endpoint,
		RedirectURL:  redirectURI,
	}
	return config
}

// Returns a URL to send to the end user for them to access as well as the state string embedded into the URL. Ensure that this state string matches the value recieved at the redirect URI.
func GetAuthCodeURL(config *auth2.Config) (string, string) {
    state := generateState()
    return config.AuthCodeURL(state, oauth2.AccessTypeOffline), state)
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

// Generate random 32 character state string
func generateState() string {
	var curr strings.Builder
	for i := 0; i < 32; i++ {
		curr.WriteString(string(stateSet[rand.Intn(len(stateSet))]))
	}
	return curr.String()
}
