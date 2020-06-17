package main

import (
	"flag"
	"fmt"
	"github.com/kelr/gundyr/auth"
	"github.com/kelr/gundyr/pubsub"
	"golang.org/x/oauth2"
)

const (
	clientID     = ""
	clientSecret = ""
	redirectURI  = "http://twitch.tv"
	userID       = ""
	tokenFile    = "token.json"
)

// Set scopes to request from the user
var scopes = []string{"channel:read:redemptions", "channel:moderate", "whispers:read"}

// Run with -a to generate a new user token
var doAuth = flag.Bool("a", false, "Generate a URL for user token authentication")

// Helper function to generate the auth code URL, generate a user credential token and flush it to a file.
func authenticate(config *oauth2.Config) error {
	// Get the URL to send to the user and the state code to protect against CSRF attacks.
	url, state := auth.GetAuthCodeURL(config)
	fmt.Println(url)
	fmt.Println("Ensure that state recieved at URI is:", state)

	// Enter the code received by the redirect URI. Ensure that the state value
	// obtained at the redirect URI matches the previous state value.
	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		fmt.Println(err)
	}

	// Obtain the user token through the code. This token can be reused as long as
	// it has not expired, but the auth code cannot be reused.
	token, err := auth.TokenExchange(config, authCode)
	if err != nil {
		fmt.Println(err)
		return err
	}
	auth.FlushTokenFile(tokenFile, token)
	return nil
}

func handleChannelPoints(event *pubsub.ChannelPointsEvent) {
	fmt.Println(event.Data.Redemption.Reward.Title)
}

func handleModActions(event *pubsub.ChatModActionsEvent) {
	fmt.Println(event.Data.ModerationAction, event.Data.TargetUserID)
}

func handleWhispers(event *pubsub.WhispersEvent) {
	fmt.Println(event.Data.Tags.Login+":", event.Data.Body, "->", event.Data.Recipient.DisplayName)
}

func main() {
	flag.Parse()

	// Setup OAuth2 configuration
	config, err := auth.NewUserAuth(clientID, clientSecret, redirectURI, &scopes)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Generate a new token if the -a flag was provided
	if *doAuth {
		authenticate(config)
		return
	}

	// Load an OAuth2 token from a file
	token, err := auth.LoadTokenFile(config, tokenFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Verify that the cached token has not expired. Refresh if so.
	newToken := auth.VerifyToken(config, token)
	if newToken.AccessToken != token.AccessToken {
		fmt.Println("Saved new token:", newToken.AccessToken)
		auth.FlushTokenFile(tokenFile, newToken)
	}
	fmt.Println("Token loaded")

	client := pubsub.NewClient(config, newToken)
	client.ListenChannelPoints(userId, handleChannelPoints)
	client.ListenChatModActions(userId, handleModActions)
	client.ListenWhispers(userId, handleWhispers)
	client.Connect()
	select {}
}
