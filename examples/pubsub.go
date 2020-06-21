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

func handleChannelPoints(event *pubsub.ChannelPointsData) {
	fmt.Println("Got Channel Points event")
	fmt.Println(event.Redemption.Reward.Title)
}

func handleModActions(event *pubsub.ChatModActionsData) {
	fmt.Println("Got Mod Action event")
	fmt.Println(event.CreatedBy, event.ModerationAction, event.Args)
}

func handleWhispers(event *pubsub.WhispersData) {
	fmt.Println("Got Whisper event")
	fmt.Println(event.Tags.Login+":", event.Body, "->", event.Recipient.DisplayName)
}

func handleSubs(event *pubsub.SubsData) {
	fmt.Println("Got Subs event")
}

func handleBits(event *pubsub.BitsData) {
	fmt.Println("Got Bits event")
}

func handleBitsBadge(event *pubsub.BitsBadgeData) {
	fmt.Println("Got Bits badge event")
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
	client := pubsub.NewClient(userID, token)

	client.ListenChannelPoints(handleChannelPoints)
	client.ListenChatModActions(handleModActions)
	client.ListenWhispers(handleWhispers)
	client.ListenSubs(handleSubs)
	client.ListenBits(handleBits)
	client.ListenBitsBadge(handleBitsBadge)

	client.Connect()
	select {}
}
