package main

import (
	"fmt"
	"log"
	"github.com/kelr/gundyr/auth"
	"github.com/kelr/gundyr/pubsub"
)

// Provide your Client ID and secret. Set your redirect URI to one that you own.
// The URI must match exactly with the one registered by your app on the Twitch Developers site
const (
	clientID     = ""
	clientSecret = ""
	redirectURI  = "http://twitch.tv"
	userID       = ""
	tokenFile    = "token.json"
)

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
	scopes := []string{"channel:read:redemptions", "channel:moderate", "whispers:read"}

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

	// Create a PubSub client and listen to the topics.
	client := pubsub.NewClient(config, token)

	client.ListenChannelPoints(userID, handleChannelPoints)
	client.ListenChatModActions(userID, handleModActions)
	client.ListenWhispers(userID, handleWhispers)
	client.Connect()
	select {}
}
